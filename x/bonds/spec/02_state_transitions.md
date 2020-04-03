# State Transitions

## Deposits

### Bond

When a deposit occurs the deposit object is affected

- remove tokens from the sending account
- add tokens to the deposit object or add them to a created deposit object
- transfer the `deposit.Amount` from the depositor's account to the `BondedPool`

### Begin Unbonding

As a part of the Unbond and Complete Unbonding state transitions Unbond
Deposit may be called.

- subtract the unbonded tokens from deposit
- add the tokens to an `UnbondingDeposit` Entry
- update the deposit or remove the deposit if there are no more tokens
- transfer the `Coins` worth of the unbonded
  tokens from the `BondedPool` to the `NotBondedPool` `ModuleAccount`

### Complete Unbonding

For undeposits which do not complete immediately, the following operations
occur when the unbonding deposit queue element matures:

- remove the entry from the `UnbondingDeposit` object
- transfer the tokens from the `NotBondedPool` `ModuleAccount` to the depositor `Account`
