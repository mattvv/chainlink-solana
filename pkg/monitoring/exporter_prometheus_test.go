package monitoring

import (
	"context"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink-solana/pkg/monitoring/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/stretchr/testify/require"
)

func TestPrometheusExporter(t *testing.T) {
	transmissionAccount := generatePublicKey()
	stateAccount := generatePublicKey()

	cfg := config.Config{}
	cfg.Feeds.Feeds = []config.Feed{
		{
			TransmissionsAccount: transmissionAccount,
			StateAccount:         stateAccount,
		},
	}

	ctx := context.Background()

	t.Run("should still publish new transmissions even if a transmitter is not set", func(t *testing.T) {
		metrics := &keepLatestMetrics{}
		exporter := NewPrometheusExporter(
			cfg.Solana, cfg.Feeds.Feeds[0],
			logger.NewNullLogger(),
			metrics,
		)

		envelope := generateTransmissionEnvelope()
		exporter.Export(ctx, envelope)
		require.Equal(t, metrics.latestTransmitter, "n/a")
	})
	t.Run("should publish a new transmission with latest transmitter", func(t *testing.T) {
		metrics := &keepLatestMetrics{}
		exporter := NewPrometheusExporter(
			cfg.Solana, cfg.Feeds.Feeds[0],
			logger.NewNullLogger(),
			metrics,
		)

		envelope1, err := generateStateEnvelope()
		require.NoError(t, err)
		exporter.Export(ctx, envelope1)

		envelope2 := generateTransmissionEnvelope()
		exporter.Export(ctx, envelope2)

		require.Equal(t, metrics.latestTransmitter, envelope1.State.Config.LatestTransmitter.String())
	})
}

type keepLatestMetrics struct {
	*devnullMetrics

	latestTransmission *big.Int
	latestTransmitter  string
}

func (k *keepLatestMetrics) SetOffchainAggregatorAnswers(answer *big.Int, contractAddress, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
	k.latestTransmission = &big.Int{}
	k.latestTransmission.Set(answer)
}

func (k *keepLatestMetrics) SetOffchainAggregatorSubmissionReceivedValues(value *big.Int, contractAddress, sender, chainID, contractStatus, contractType, feedName, feedPath, networkID, networkName string) {
	k.latestTransmission = &big.Int{}
	k.latestTransmission.Set(value)
	k.latestTransmitter = sender
}