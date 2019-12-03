package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/cosmos/cosmos-sdk/store"

	"github.com/cosmos/cosmos-sdk/baseapp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"

	xar "github.com/xar-network/xar-network/app"
)

func runHackCmd(cmd *cobra.Command, args []string) error {

	if len(args) != 1 {
		return fmt.Errorf("expected 1 arg")
	}

	// ".xard"
	dataDir := args[0]
	dataDir = path.Join(dataDir, "data")

	// load the app
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	db, err := sdk.NewLevelDB("xar", dataDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	app, keyMain, keyStaking, stakingKeeper := xar.NewXarAppUNSAFE(
		logger, db, nil, false, 0, baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))))

	// print some info
	id := app.LastCommitID()
	lastBlockHeight := app.LastBlockHeight()
	fmt.Println("ID", id)
	fmt.Println("LastBlockHeight", lastBlockHeight)

	//----------------------------------------------------
	// XXX: start hacking!
	//----------------------------------------------------
	// eg. xar-6001 testnet bug
	// We paniced when iterating through the "bypower" keys.
	// The following powerKey was there, but the corresponding "trouble" validator did not exist.
	// So here we do a binary search on the past states to find when the powerKey first showed up ...

	// operator of the validator the bonds, gets jailed, later unbonds, and then later is still found in the bypower store
	trouble := hexToBytes("D3DC0FF59F7C3B548B7AFA365561B87FD0208AF8")
	// this is his "bypower" key
	powerKey := hexToBytes("05303030303030303030303033FFFFFFFFFFFF4C0C0000FFFED3DC0FF59F7C3B548B7AFA365561B87FD0208AF8")

	topHeight := lastBlockHeight
	bottomHeight := int64(0)
	checkHeight := topHeight
	for {
		// load the given version of the state
		err = app.LoadVersion(checkHeight, keyMain)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		ctx := app.NewContext(true, abci.Header{})

		// check for the powerkey and the validator from the store
		store := ctx.KVStore(keyStaking)
		res := store.Get(powerKey)
		val, _ := stakingKeeper.GetValidator(ctx, trouble)
		fmt.Println("checking height", checkHeight, res, val)
		if res == nil {
			bottomHeight = checkHeight
		} else {
			topHeight = checkHeight
		}
		checkHeight = (topHeight + bottomHeight) / 2
	}
}

func hexToBytes(h string) []byte {
	trouble, err := hex.DecodeString(h)
	if err != nil {
		return nil
	}
	return trouble

}