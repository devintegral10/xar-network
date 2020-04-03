# Governors module specification

## Abstract

Governors module allows to bond an amount of XAR tokens to gain a right to
vote on proposals in `gov` module. The weight of a governor is the amount
of bonded tokens.

The mechanism of bonding is implemented via `bonds` keeper. It's made similar to
bonding delegators in `staking` module, but without validators.

## Contents

1. **[Messages](01_messages.md)**
2. **[Events](02_events.md)**
3. **[Parameters](03_params.md)**

See `Bonds` spec for information about inner state, state transitions, hooks.
