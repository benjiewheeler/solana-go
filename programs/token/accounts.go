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

package token

import (
	"encoding/binary"

	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
)

type Mint struct {
	// Optional authority used to mint new tokens. The mint authority may only be provided during
	// mint creation. If no mint authority is present then the mint has a fixed supply and no
	// further tokens may be minted.
	MintAuthority *ag_solanago.PublicKey `bin:"optional"`

	// Total supply of tokens.
	Supply uint64

	// Number of base 10 digits to the right of the decimal place.
	Decimals uint8

	// Is `true` if this structure has been initialized
	IsInitialized bool

	// Optional authority to freeze token accounts.
	FreezeAuthority *ag_solanago.PublicKey `bin:"optional"`
}

func (mint *Mint) UnmarshalWithDecoder(dec *ag_binary.Decoder) (err error) {
	{
		v, err := dec.ReadUint32(binary.LittleEndian)
		if err != nil {
			return err
		}
		if v == 1 {
			v, err := dec.ReadNBytes(32)
			if err != nil {
				return err
			}
			mint.MintAuthority = ag_solanago.PublicKeyFromBytes(v).ToPointer()
		} else {
			// discard:
			_, err := dec.ReadNBytes(32)
			if err != nil {
				return err
			}
		}
	}
	{
		v, err := dec.ReadUint64(binary.LittleEndian)
		if err != nil {
			return err
		}
		mint.Supply = v
	}
	{
		v, err := dec.ReadUint8()
		if err != nil {
			return err
		}
		mint.Decimals = v
	}
	{
		v, err := dec.ReadBool()
		if err != nil {
			return err
		}
		mint.IsInitialized = v
	}
	{
		v, err := dec.ReadUint32(binary.LittleEndian)
		if err != nil {
			return err
		}
		if v == 1 {
			v, err := dec.ReadNBytes(32)
			if err != nil {
				return err
			}
			mint.FreezeAuthority = ag_solanago.PublicKeyFromBytes(v).ToPointer()
		} else {
			// discard:
			_, err := dec.ReadNBytes(32)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (mint Mint) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	{
		if mint.MintAuthority == nil {
			err = encoder.WriteUint32(0, binary.LittleEndian)
			if err != nil {
				return err
			}
			empty := ag_solanago.PublicKey{}
			err = encoder.WriteBytes(empty[:], false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteUint32(1, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteBytes(mint.MintAuthority[:], false)
			if err != nil {
				return err
			}
		}
	}
	err = encoder.WriteUint64(mint.Supply, binary.LittleEndian)
	if err != nil {
		return err
	}
	err = encoder.WriteUint8(mint.Decimals)
	if err != nil {
		return err
	}
	err = encoder.WriteBool(mint.IsInitialized)
	if err != nil {
		return err
	}
	{
		if mint.FreezeAuthority == nil {
			err = encoder.WriteUint32(0, binary.LittleEndian)
			if err != nil {
				return err
			}
			empty := ag_solanago.PublicKey{}
			err = encoder.WriteBytes(empty[:], false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteUint32(1, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteBytes(mint.FreezeAuthority[:], false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type Account struct {
	// The mint associated with this account
	Mint ag_solanago.PublicKey

	// The owner of this account.
	Owner ag_solanago.PublicKey

	// The amount of tokens this account holds.
	Amount uint64

	// If `delegate` is `Some` then `delegated_amount` represents
	// the amount authorized by the delegate
	Delegate *ag_solanago.PublicKey `bin:"optional"`

	// The account's state
	State AccountState

	// If is_some, this is a native token, and the value logs the rent-exempt reserve. An Account
	// is required to be rent-exempt, so the value is used by the Processor to ensure that wrapped
	// SOL accounts do not drop below this threshold.
	IsNative *uint64 `bin:"optional"`

	// The amount delegated
	DelegatedAmount uint64

	// Optional authority to close the account.
	CloseAuthority *ag_solanago.PublicKey `bin:"optional"`
}

func (mint *Account) UnmarshalWithDecoder(dec *ag_binary.Decoder) (err error) {
	{
		v, err := dec.ReadNBytes(32)
		if err != nil {
			return err
		}
		mint.Mint = ag_solanago.PublicKeyFromBytes(v)
	}
	{
		v, err := dec.ReadNBytes(32)
		if err != nil {
			return err
		}
		mint.Owner = ag_solanago.PublicKeyFromBytes(v)
	}
	{
		v, err := dec.ReadUint64(binary.LittleEndian)
		if err != nil {
			return err
		}
		mint.Amount = v
	}
	{
		v, err := dec.ReadUint32(binary.LittleEndian)
		if err != nil {
			return err
		}
		if v == 1 {
			v, err := dec.ReadNBytes(32)
			if err != nil {
				return err
			}
			mint.Delegate = ag_solanago.PublicKeyFromBytes(v).ToPointer()
		} else {
			// discard:
			_, err := dec.ReadNBytes(32)
			if err != nil {
				return err
			}
		}
	}
	{
		v, err := dec.ReadUint8()
		if err != nil {
			return err
		}
		mint.State = AccountState(v)
	}
	{
		v, err := dec.ReadUint32(binary.LittleEndian)
		if err != nil {
			return err
		}
		if v == 1 {
			v, err := dec.ReadUint64(ag_binary.LE)
			if err != nil {
				return err
			}
			mint.IsNative = &v
		} else {
			// discard:
			_, err := dec.ReadUint64(ag_binary.LE)
			if err != nil {
				return err
			}
		}
	}
	{
		v, err := dec.ReadUint64(binary.LittleEndian)
		if err != nil {
			return err
		}
		mint.DelegatedAmount = v
	}
	{
		v, err := dec.ReadUint32(binary.LittleEndian)
		if err != nil {
			return err
		}
		if v == 1 {
			v, err := dec.ReadNBytes(32)
			if err != nil {
				return err
			}
			mint.CloseAuthority = ag_solanago.PublicKeyFromBytes(v).ToPointer()
		} else {
			// discard:
			_, err := dec.ReadNBytes(32)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (mint Account) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	{
		err = encoder.WriteBytes(mint.Mint[:], false)
		if err != nil {
			return err
		}
	}
	{
		err = encoder.WriteBytes(mint.Owner[:], false)
		if err != nil {
			return err
		}
	}
	{
		err = encoder.WriteUint64(mint.Amount, ag_binary.LE)
		if err != nil {
			return err
		}
	}
	{
		if mint.Delegate == nil {
			err = encoder.WriteUint32(0, binary.LittleEndian)
			if err != nil {
				return err
			}
			empty := ag_solanago.PublicKey{}
			err = encoder.WriteBytes(empty[:], false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteUint32(1, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteBytes(mint.Delegate[:], false)
			if err != nil {
				return err
			}
		}
	}
	err = encoder.WriteUint8(uint8(mint.State))
	if err != nil {
		return err
	}
	{
		if mint.IsNative == nil {
			err = encoder.WriteUint32(0, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteUint64(0, ag_binary.LE)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteUint32(1, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteUint64(*mint.IsNative, ag_binary.LE)
			if err != nil {
				return err
			}
		}
	}
	{
		err = encoder.WriteUint64(mint.DelegatedAmount, ag_binary.LE)
		if err != nil {
			return err
		}
	}
	{
		if mint.CloseAuthority == nil {
			err = encoder.WriteUint32(0, binary.LittleEndian)
			if err != nil {
				return err
			}
			empty := ag_solanago.PublicKey{}
			err = encoder.WriteBytes(empty[:], false)
			if err != nil {
				return err
			}
		} else {
			err = encoder.WriteUint32(1, binary.LittleEndian)
			if err != nil {
				return err
			}
			err = encoder.WriteBytes(mint.CloseAuthority[:], false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type Multisig struct {
	// Number of signers required
	M uint8
	// Number of valid signers
	N uint8
	// Is `true` if this structure has been initialized
	IsInitialized bool
	// Signer public keys
	Signers [MAX_SIGNERS]ag_solanago.PublicKey
}

func (obj *Multisig) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	if err = encoder.Encode(&obj.M); err != nil {
		return err
	}
	if err = encoder.Encode(&obj.N); err != nil {
		return err
	}
	if err = encoder.Encode(&obj.IsInitialized); err != nil {
		return err
	}
	if err = encoder.Encode(&obj.Signers); err != nil {
		return err
	}
	return nil
}

func (obj *Multisig) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	if err = decoder.Decode(&obj.M); err != nil {
		return err
	}
	if err = decoder.Decode(&obj.N); err != nil {
		return err
	}
	if err = decoder.Decode(&obj.IsInitialized); err != nil {
		return err
	}
	if err = decoder.Decode(&obj.Signers); err != nil {
		return err
	}
	return nil
}
