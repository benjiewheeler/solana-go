// Copyright 2021 github.com/gagliardetto
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package associatedtokenaccount

import (
	"fmt"

	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

type Create struct {
	// [0] = [WRITE, SIGNER] Payer
	// ··········· Funding account
	//
	// [1] = [WRITE] AssociatedTokenAccount
	// ··········· Associated token account address to be created
	//
	// [2] = [] Wallet
	// ··········· Wallet address for the new associated token account
	//
	// [3] = [] TokenMint
	// ··········· The token mint for the new associated token account
	//
	// [4] = [] SystemProgram
	// ··········· System program ID
	//
	// [5] = [] TokenProgram
	// ··········· SPL token program ID
	//
	// [6] = [] SysVarRent
	// ··········· SysVarRentPubkey
	ag_solanago.AccountMetaSlice `bin:"-" borsh_skip:"true"`
}

// NewCreateInstructionBuilder creates a new `Create` instruction builder.
func NewCreateInstructionBuilder() *Create {
	nd := &Create{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 7),
	}
	nd.AccountMetaSlice[4] = ag_solanago.Meta(ag_solanago.SystemProgramID)
	nd.AccountMetaSlice[5] = ag_solanago.Meta(ag_solanago.TokenProgramID)
	nd.AccountMetaSlice[6] = ag_solanago.Meta(ag_solanago.SysVarRentPubkey)
	return nd
}

func (inst *Create) SetPayer(payer ag_solanago.PublicKey) *Create {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(payer).WRITE().SIGNER()
	return inst
}

func (inst Create) GetPayer() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[0]
}

func (inst *Create) SetAssociatedTokenAccount(associatedTokenAccount ag_solanago.PublicKey) *Create {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(associatedTokenAccount).WRITE()
	return inst
}

func (inst Create) GetAssociatedTokenAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[1]
}

func (inst *Create) SetWallet(wallet ag_solanago.PublicKey) *Create {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(wallet)
	return inst
}

func (inst Create) GetWallet() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[2]
}

func (inst *Create) SetMint(mint ag_solanago.PublicKey) *Create {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(mint)
	return inst
}

func (inst Create) GetMint() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[3]
}

func (inst *Create) SetTokenProgramID(tokenProgramID ag_solanago.PublicKey) *Create {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(tokenProgramID)
	return inst
}

func (inst Create) GetTokenProgramID() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice[5]
}

func (inst Create) Build() *Instruction {
	if ata := inst.AccountMetaSlice[1]; ata == nil || ata.PublicKey.IsZero() {
		// Find the associatedTokenAddress;
		associatedTokenAddress, _, _ := ag_solanago.FindAssociatedTokenAddress(
			inst.GetWallet().PublicKey,
			inst.GetMint().PublicKey,
			inst.GetTokenProgramID().PublicKey,
		)

		inst.SetAssociatedTokenAccount(associatedTokenAddress)
	}

	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: ag_binary.TypeIDFromUint8(Instruction_Create),
	}}
}

// ValidateAndBuild validates the instruction accounts.
// If there is a validation error, return the error.
// Otherwise, build and return the instruction.
func (inst Create) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *Create) Validate() error {
	if ata := inst.AccountMetaSlice[1]; ata == nil || ata.PublicKey.IsZero() {
		// Find the associatedTokenAddress;
		associatedTokenAddress, _, _ := ag_solanago.FindAssociatedTokenAddress(
			inst.GetWallet().PublicKey,
			inst.GetMint().PublicKey,
			inst.GetTokenProgramID().PublicKey,
		)

		inst.SetAssociatedTokenAccount(associatedTokenAddress)
	}

	// Check whether all accounts are set:
	for accIndex, acc := range inst.AccountMetaSlice {
		if acc == nil {
			return fmt.Errorf("ins.AccountMetaSlice[%v] is not set", accIndex)
		}
	}
	return nil
}

func (inst *Create) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("Create")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {
					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=7").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("                 payer", inst.Get(0)))
						accountsBranch.Child(ag_format.Meta("associatedTokenAddress", inst.Get(1)))
						accountsBranch.Child(ag_format.Meta("                wallet", inst.Get(2)))
						accountsBranch.Child(ag_format.Meta("             tokenMint", inst.Get(3)))
						accountsBranch.Child(ag_format.Meta("         systemProgram", inst.Get(4)))
						accountsBranch.Child(ag_format.Meta("          tokenProgram", inst.Get(5)))
						accountsBranch.Child(ag_format.Meta("            sysVarRent", inst.Get(6)))
					})
				})
		})
}

func (inst Create) MarshalWithEncoder(encoder *ag_binary.Encoder) error {
	return encoder.WriteBytes([]byte{}, false)
}

func (inst *Create) UnmarshalWithDecoder(decoder *ag_binary.Decoder) error {
	return nil
}

func NewCreateInstruction(
	payer ag_solanago.PublicKey,
	walletAddress ag_solanago.PublicKey,
	splTokenMintAddress ag_solanago.PublicKey,
	tokenProgramID ag_solanago.PublicKey,
) *Create {
	return NewCreateInstructionBuilder().
		SetPayer(payer).
		SetWallet(walletAddress).
		SetMint(splTokenMintAddress).
		SetTokenProgramID(tokenProgramID)
}
