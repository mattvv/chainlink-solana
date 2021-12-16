// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package ocr_2

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// SetRequesterAccessController is the `setRequesterAccessController` instruction.
type SetRequesterAccessController struct {

	// [0] = [WRITE] state
	//
	// [1] = [SIGNER] authority
	//
	// [2] = [] accessController
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewSetRequesterAccessControllerInstructionBuilder creates a new `SetRequesterAccessController` instruction builder.
func NewSetRequesterAccessControllerInstructionBuilder() *SetRequesterAccessController {
	nd := &SetRequesterAccessController{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 3),
	}
	return nd
}

// SetStateAccount sets the "state" account.
func (inst *SetRequesterAccessController) SetStateAccount(state ag_solanago.PublicKey) *SetRequesterAccessController {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(state).WRITE()
	return inst
}

// GetStateAccount gets the "state" account.
func (inst *SetRequesterAccessController) GetStateAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// SetAuthorityAccount sets the "authority" account.
func (inst *SetRequesterAccessController) SetAuthorityAccount(authority ag_solanago.PublicKey) *SetRequesterAccessController {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(authority).SIGNER()
	return inst
}

// GetAuthorityAccount gets the "authority" account.
func (inst *SetRequesterAccessController) GetAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[1]
}

// SetAccessControllerAccount sets the "accessController" account.
func (inst *SetRequesterAccessController) SetAccessControllerAccount(accessController ag_solanago.PublicKey) *SetRequesterAccessController {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(accessController)
	return inst
}

// GetAccessControllerAccount gets the "accessController" account.
func (inst *SetRequesterAccessController) GetAccessControllerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[2]
}

func (inst SetRequesterAccessController) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_SetRequesterAccessController,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst SetRequesterAccessController) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *SetRequesterAccessController) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.State is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Authority is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.AccessController is not set")
		}
	}
	return nil
}

func (inst *SetRequesterAccessController) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("SetRequesterAccessController")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=3]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("           state", inst.AccountMetaSlice[0]))
						accountsBranch.Child(ag_format.Meta("       authority", inst.AccountMetaSlice[1]))
						accountsBranch.Child(ag_format.Meta("accessController", inst.AccountMetaSlice[2]))
					})
				})
		})
}

func (obj SetRequesterAccessController) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *SetRequesterAccessController) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewSetRequesterAccessControllerInstruction declares a new SetRequesterAccessController instruction with the provided parameters and accounts.
func NewSetRequesterAccessControllerInstruction(
	// Accounts:
	state ag_solanago.PublicKey,
	authority ag_solanago.PublicKey,
	accessController ag_solanago.PublicKey) *SetRequesterAccessController {
	return NewSetRequesterAccessControllerInstructionBuilder().
		SetStateAccount(state).
		SetAuthorityAccount(authority).
		SetAccessControllerAccount(accessController)
}