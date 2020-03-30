package simulation

import (
	"bytes"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/xar-network/xar-network/x/bonds"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding type
func DecodeStore(cdc *codec.Codec, kvA, kvB cmn.KVPair) string {
	switch {

	case bytes.Equal(kvA.Key[:1], bonds.DepositKey):
		var depositA, depositB bonds.Deposit
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &depositA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &depositB)
		return fmt.Sprintf("%v\n%v", depositA, depositB)

	case bytes.Equal(kvA.Key[:1], bonds.UnbondingDepositKey):
		var ubdA, ubdB bonds.UnbondingDeposit
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &ubdA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &ubdB)
		return fmt.Sprintf("%v\n%v", ubdA, ubdB)

	default:
		panic(fmt.Sprintf("invalid governors key prefix %X", kvA.Key[:1]))
	}
}
