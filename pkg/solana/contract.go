package solana

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"

	"golang.org/x/sync/singleflight"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
)

var (
	configVersion uint8 = 1
	rpcCommitment       = rpc.CommitmentConfirmed

	// error declarations
	errCursorLength       = errors.New("incorrect cursor length")
	errTransmissionLength = errors.New("incorrect transmission length")
)

type ContractTracker struct {
	// on-chain program + 2x state accounts (state + transmissions)
	ProgramID       solana.PublicKey
	StateID         solana.PublicKey
	TransmissionsID solana.PublicKey
	StoreProgramID  solana.PublicKey

	// private key for the transmission signing
	Transmitter TransmissionSigner

	// tracked contract state
	state  State
	store  solana.PublicKey
	answer Answer

	// dependencies
	client *Client
	lggr   Logger

	// provides a duplicate function call suppression mechanism
	requestGroup *singleflight.Group
}

func NewTracker(spec OCR2Spec, client *Client, transmitter TransmissionSigner, lggr Logger) ContractTracker {
	return ContractTracker{
		ProgramID:       spec.ProgramID,
		StateID:         spec.StateID,
		StoreProgramID:  spec.StoreProgramID,
		TransmissionsID: spec.TransmissionsID,
		Transmitter:     transmitter,
		client:          client,
		lggr:            lggr,
		requestGroup:    &singleflight.Group{},
	}
}

// fetch + decode + store raw state
func (c *ContractTracker) fetchState(ctx context.Context) error {
	c.lggr.Debugf("fetch state for account: %s", c.StateID.String())

	// make single flight request
	v, err, shared := c.requestGroup.Do("state", func() (interface{}, error) {
		state, _, err := GetState(ctx, c.client.rpc, c.StateID)
		return state, err
	})

	if err != nil {
		return err
	}

	c.state = v.(State)
	c.lggr.Debugf("state fetched for account: %s, shared: %t, result (config digest): %v", c.StateID, shared, hex.EncodeToString(c.state.Config.LatestConfigDigest[:]))

	// Fetch the store address associated to the feed
	offset := uint64(8 + 1) // Discriminator (8 bytes) + Version (u8)
	length := uint64(solana.PublicKeyLength)
	res, err := c.client.rpc.GetAccountInfoWithOpts(ctx, c.state.Transmissions, &rpc.GetAccountInfoOpts{
		Encoding:   "base64",
		Commitment: rpcCommitment,
		DataSlice: &rpc.DataSlice{
			Offset: &offset,
			Length: &length,
		},
	})
	if err != nil {
		return err
	}
	c.store = solana.PublicKeyFromBytes(res.Value.Data.GetBinary())
	c.lggr.Debugf("store fetched for feed: %s", c.store)

	return nil
}

func (c *ContractTracker) fetchLatestTransmission(ctx context.Context) error {
	c.lggr.Debugf("fetch latest transmission for account: %s", c.TransmissionsID)

	// make single flight request
	v, err, shared := c.requestGroup.Do("transmissions.latest", func() (interface{}, error) {
		answer, _, err := GetLatestTransmission(ctx, c.client.rpc, c.TransmissionsID)
		return answer, err
	})

	if err != nil {
		return err
	}

	c.lggr.Debugf("latest transmission fetched for account: %s, shared: %t, result: %v", c.TransmissionsID, shared, v)

	c.answer = v.(Answer)
	return nil
}

func GetState(ctx context.Context, client *rpc.Client, account solana.PublicKey) (State, uint64, error) {
	res, err := client.GetAccountInfoWithOpts(ctx, account, &rpc.GetAccountInfoOpts{
		Encoding:   "base64",
		Commitment: rpcCommitment,
	})
	if err != nil {
		return State{}, 0, fmt.Errorf("failed to fetch state account at address '%s': %w", account.String(), err)
	}

	var state State
	if err := bin.NewBinDecoder(res.Value.Data.GetBinary()).Decode(&state); err != nil {
		return State{}, 0, fmt.Errorf("failed to decode state account data: %w", err)
	}

	// validation for config version
	if configVersion != state.Version {
		return State{}, 0, fmt.Errorf("decoded config version (%d) does not match expected config version (%d)", state.Version, configVersion)
	}

	blockNum := res.RPCContext.Context.Slot
	return state, blockNum, nil
}

func GetLatestTransmission(ctx context.Context, client *rpc.Client, account solana.PublicKey) (Answer, uint64, error) {
	offset := CursorOffset
	length := CursorLen * 2
	transmissionLen := TransmissionLen

	// query for cursor
	res, err := client.GetAccountInfoWithOpts(ctx, account, &rpc.GetAccountInfoOpts{
		Encoding:   "base64",
		Commitment: rpcCommitment,
		DataSlice: &rpc.DataSlice{
			Offset: &offset,
			Length: &length,
		},
	})
	if err != nil {
		return Answer{}, 0, errors.Wrap(err, "error on rpc.GetAccountInfo [cursor]")
	}

	// parse little endian length & cursor
	c := res.Value.Data.GetBinary()
	buf := bytes.NewReader(c)

	var cursor uint32
	var liveLength uint32
	err = binary.Read(buf, binary.LittleEndian, &liveLength)
	if err != nil {
		return Answer{}, 0, errCursorLength
	}
	err = binary.Read(buf, binary.LittleEndian, &cursor)
	if err != nil {
		return Answer{}, 0, errCursorLength
	}

	if cursor == 0 { // handle array wrap
		cursor = liveLength
	}
	cursor-- // cursor indicates index for new answer, latest answer is in previous index

	// fetch transmission
	var transmissionOffset uint64 = 8 + 128 + (uint64(cursor) * transmissionLen)

	res, err = client.GetAccountInfoWithOpts(ctx, account, &rpc.GetAccountInfoOpts{
		Encoding:   "base64",
		Commitment: rpcCommitment,
		DataSlice: &rpc.DataSlice{
			Offset: &transmissionOffset,
			Length: &transmissionLen,
		},
	})
	if err != nil {
		return Answer{}, 0, errors.Wrap(err, "error on rpc.GetAccountInfo [transmission]")
	}

	t := res.Value.Data.GetBinary()
	if len(t) != int(transmissionLen) { // validate length
		return Answer{}, 0, errTransmissionLength
	}

	var timestamp uint64
	raw := make([]byte, 16)

	buf = bytes.NewReader(t)
	err = binary.Read(buf, binary.LittleEndian, &timestamp)
	if err != nil {
		return Answer{}, 0, err
	}

	// TODO: we could use ag_binary.Int128 instead
	err = binary.Read(buf, binary.LittleEndian, &raw)
	if err != nil {
		return Answer{}, 0, err
	}
	// reverse slice to change from little endian to big endian
	for i, j := 0, len(raw)-1; i < j; i, j = i+1, j-1 {
		raw[i], raw[j] = raw[j], raw[i]
	}

	return Answer{
		Data:      big.NewInt(0).SetBytes(raw[:]),
		Timestamp: timestamp,
	}, res.RPCContext.Context.Slot, nil
}
