# Events

The governors module emits the following events:

## EndBlocker

| Type                  | Attribute Key         | Attribute Value       |
|-----------------------|-----------------------|-----------------------|
| complete_unbonding    | depositor             | {depositorAddress}    |

## Handlers

### MsgBond

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| bond | validator     | {validatorAddress} |
| bond | amount        | {depositAmount} |
| message  | module        | governors            |
| message  | action        | bond           |
| message  | sender        | {senderAddress}    |

### MsgUnbond

| Type    | Attribute Key       | Attribute Value    |
|---------|---------------------|--------------------|
| unbond  | amount              | {unbondAmount}     |
| unbond  | completion_time [0] | {completionTime}   |
| message | module              | governors            |
| message | action              | begin_unbonding    |
| message | sender              | {senderAddress}    |

* [0] Time is formatted in the RFC3339 standard
