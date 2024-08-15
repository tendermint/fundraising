package sample

import (
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// PubKey returns a sample public key.
func PubKey() cryptotypes.PubKey {
	return ed25519.GenPrivKey().PubKey()
}

// AccAddress returns a sample account address.
func AccAddress() sdk.AccAddress {
	addr := PubKey().Address()
	return sdk.AccAddress(addr)
}

// Address returns a sample account address string.
func Address() string {
	return AccAddress().String()
}
