// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package ocr_2

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// PayOracles is the `payOracles` instruction.
type PayOracles struct {

	// [0] = [WRITE] state
	//
	// [1] = [SIGNER] authority
	//
	// [2] = [] accessController
	//
	// [3] = [WRITE] tokenVault
	//
	// [4] = [] vaultAuthority
	//
	// [5] = [] tokenProgram
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewPayOraclesInstructionBuilder creates a new `PayOracles` instruction builder.
func NewPayOraclesInstructionBuilder() *PayOracles {
	nd := &PayOracles{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 6),
	}
	return nd
}

// SetStateAccount sets the "state" account.
func (inst *PayOracles) SetStateAccount(state ag_solanago.PublicKey) *PayOracles {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(state).WRITE()
	return inst
}

// GetStateAccount gets the "state" account.
func (inst *PayOracles) GetStateAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

// SetAuthorityAccount sets the "authority" account.
func (inst *PayOracles) SetAuthorityAccount(authority ag_solanago.PublicKey) *PayOracles {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(authority).SIGNER()
	return inst
}

// GetAuthorityAccount gets the "authority" account.
func (inst *PayOracles) GetAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[1]
}

// SetAccessControllerAccount sets the "accessController" account.
func (inst *PayOracles) SetAccessControllerAccount(accessController ag_solanago.PublicKey) *PayOracles {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(accessController)
	return inst
}

// GetAccessControllerAccount gets the "accessController" account.
func (inst *PayOracles) GetAccessControllerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[2]
}

// SetTokenVaultAccount sets the "tokenVault" account.
func (inst *PayOracles) SetTokenVaultAccount(tokenVault ag_solanago.PublicKey) *PayOracles {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(tokenVault).WRITE()
	return inst
}

// GetTokenVaultAccount gets the "tokenVault" account.
func (inst *PayOracles) GetTokenVaultAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[3]
}

// SetVaultAuthorityAccount sets the "vaultAuthority" account.
func (inst *PayOracles) SetVaultAuthorityAccount(vaultAuthority ag_solanago.PublicKey) *PayOracles {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(vaultAuthority)
	return inst
}

// GetVaultAuthorityAccount gets the "vaultAuthority" account.
func (inst *PayOracles) GetVaultAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[4]
}

// SetTokenProgramAccount sets the "tokenProgram" account.
func (inst *PayOracles) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *PayOracles {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "tokenProgram" account.
func (inst *PayOracles) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[5]
}

func (inst PayOracles) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_PayOracles,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst PayOracles) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *PayOracles) Validate() error {
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
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.TokenVault is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.VaultAuthority is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
	}
	return nil
}

func (inst *PayOracles) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("PayOracles")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=6]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("           state", inst.AccountMetaSlice[0]))
						accountsBranch.Child(ag_format.Meta("       authority", inst.AccountMetaSlice[1]))
						accountsBranch.Child(ag_format.Meta("accessController", inst.AccountMetaSlice[2]))
						accountsBranch.Child(ag_format.Meta("      tokenVault", inst.AccountMetaSlice[3]))
						accountsBranch.Child(ag_format.Meta("  vaultAuthority", inst.AccountMetaSlice[4]))
						accountsBranch.Child(ag_format.Meta("    tokenProgram", inst.AccountMetaSlice[5]))
					})
				})
		})
}

func (obj PayOracles) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *PayOracles) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewPayOraclesInstruction declares a new PayOracles instruction with the provided parameters and accounts.
func NewPayOraclesInstruction(
	// Accounts:
	state ag_solanago.PublicKey,
	authority ag_solanago.PublicKey,
	accessController ag_solanago.PublicKey,
	tokenVault ag_solanago.PublicKey,
	vaultAuthority ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey) *PayOracles {
	return NewPayOraclesInstructionBuilder().
		SetStateAccount(state).
		SetAuthorityAccount(authority).
		SetAccessControllerAccount(accessController).
		SetTokenVaultAccount(tokenVault).
		SetVaultAuthorityAccount(vaultAuthority).
		SetTokenProgramAccount(tokenProgram)
}