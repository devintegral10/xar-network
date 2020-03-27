package cli

import (
	flag "github.com/spf13/pflag"
)

// nolint
const (
	FlagAmount = "amount"
)

// common flagsets to add to various functions
var (
	FsAmount = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	FsAmount.String(FlagAmount, "", "Amount of coins to bond")
}
