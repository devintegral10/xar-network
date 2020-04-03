# Messages

In this section we describe the processing of the governors messages and
the corresponding updates to the state. All created/modified state objects
specified by each message are defined within the state section of `bonds` spec.

## MsgBond

Within this message the depositor provides coins, and in return receives
voting weight in `gov` module.

```go
type MsgBond struct {
  DepositorAddr sdk.AccAddress
  Amount        sdk.Coin
}
```

This message is expected to fail if:

- the `Amount` `Coin` has a denomination different than one defined by `params.BondDenom`
- the `Amount` is zero.

If an existing `Deposit` object for provided addresses does not already
exist then it is created as part of this message otherwise the existing
`Deposit` is updated to include the newly received shares.

## MsgBeginUnbonding

The begin unbonding message allows depositors to unbond their tokens.

```go
type MsgBeginUnbonding struct {
  DepositorAddr sdk.AccAddress
  Amount        sdk.Coin
}
```

This message is expected to fail if:

- the deposit doesn't exist
- the deposit has less tokens than `Amount`
- existing `UnbondingDeposit` has maximum entries as defined by `params.MaxEntries`
- the `Amount` has a denomination different than one defined by `params.BondDenom`
- the `Amount` is zero

When this message is processed the following actions occur:

- the deposit's `Tokens` are reduced by the message `Amount`
- with those removed tokens:
  - add them to an entry in `UnbondingDeposit` (create `UnbondingDeposit` if it doesn't exist) with a completion time a full unbonding period from the current time. Update pool shares to reduce BondedTokens and increase NotBondedTokens by token worth of the shares.
- if there are no more `Tokens` in the deposit, then the deposit object is removed from the store
