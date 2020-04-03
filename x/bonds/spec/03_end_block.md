# End-Block

Each abci end block call, the operations to update queues are specified to execute.

## Queues

Certain state-transitions are not instantaneous but take place
over a duration of time (typically the unbonding period). When these
transitions are mature certain operations must take place in order to complete
the state operation. This is achieved through the use of queues which are
checked/processed at the end of each block.

### Unbonding Deposits

Complete the unbonding of all mature `UnbondingDeposits.Entries` within the
`UnbondingDeposits` queue with the following procedure:

- transfer the balance coins to the depositor's wallet address
- remove the mature entry from `UnbondingDeposit.Entries`
- remove the `UnbondingDeposit` object from the store if there are no
  remaining entries.
