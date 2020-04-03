# Bonds library specification

## Abstract

Bonds isn't a standalone module, but a library keeper which may be
composed/aggregated into modules.

It provides basic functionality for bonding/unbonding any deposits, similar to
delegators from `staking` module but without validators.

Bond denom, maximum number of pending unbonding requests may be configured via params.

Similar to `staking` module, `bonds` keeper tracks bonded/unbonding supply.

## Contents

1. **[State](01_state.md)**
2. **[State Transitions](02_state_transitions.md)**
3. **[End-Block ](03_end_block.md)**
4. **[Hooks](04_hooks.md)**
