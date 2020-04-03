# State

## Params

Params is a module-wide configuration structure that stores system parameters
and defines overall functioning of the bonds keeper.

- Params: `Paramsspace("bonds") -> amino(params)`

```go
type Params struct {
	UnbondingTime time.Duration // time duration of unbonding
	MaxEntries    uint16        // max entries for either unbonding deposit or redeposit (per pair/trio)
	
	BondDenom         string    // bondable coin denomination
	BondedPoolName    string    // pool for bonded tokens
	NotBondedPoolName string    // pool for not bonded tokens
}
```

## Deposit

Deposits are identified by `DepositorAddr` (the address of the depositor).

Deposits are indexed in the store as follows:

- Deposit: `0x21 | DepositorAddr -> amino(deposit)`

Holders may bond coins; under this circumstance their
funds are held in a `Deposit` data structure. It is owned by one
depositor. The sender of the transaction is the owner of the bond.

```go
type Deposit struct {
    DepositorAddr sdk.AccAddress
    Tokens        sdk.Int
}
```

## UnbondingDeposit

Tokens in a `Deposit` can be unbonded, but they must for some time exist as
an `UnbondingDeposit`.

`UnbondingDeposit` are indexed in the store as:

- UnbondingDeposit: `0x22 | DepositorAddr ->
   amino(unbondingDeposit)`

The map here is used in queries, to lookup all unbonding deposits for
a given depositor.

A UnbondingDeposit object is created every time an unbonding is initiated.

```go
type UnbondingDeposit struct {
    DepositorAddr sdk.AccAddress          // depositor
    Entries       []UnbondingDepositEntry // unbonding deposit entries
}
```

## Queues

All queues objects are sorted by timestamp. The time used within any queue is
first rounded to the nearest nanosecond then sorted. The sortable time format
used is a slight modification of the RFC3339Nano and uses the the format string
`"2006-01-02T15:04:05.000000000"`. Notably this format:

- right pads all zeros
- drops the time zone info (uses UTC)

In all cases, the stored timestamp represents the maturation time of the queue
element.

### UnbondingDepositQueue

For the purpose of tracking progress of unbonding delegations the unbonding
delegations queue is kept.

- UnbondingDeposit: `0x31 | format(time) -> DepositorAddr`
