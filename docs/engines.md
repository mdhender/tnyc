# Engines

The game is processed by multiple engines. Each engine has one stage of the
turn-processing workflow and communicates with other stages through persisted
data rather than combining all processing into one component.

## Orders preprocessor

The orders preprocessor:

1. determines which actor submitted each order;
2. verifies that the actor is allowed to issue it;
3. validates the order; and
4. records whether the order is valid.

Validation does not update game state. It produces the validated orders that
the executor consumes.

## Orders executor

The orders executor reads validated orders and executes them. Execution updates
the fact tables that hold game state. Invalid or unvalidated orders are not
executed.

## Reports engine

The reports engine reads the fact tables after execution and produces reports
for players and administrators. Reports are derived output; generating them
does not change game state.

## Processing flow

```text
Actors
  │
  ▼
submitted orders
  │
  ▼
Orders preprocessor ── invalid ──▶ rejected orders
  │
  │ validated
  ▼
Orders executor
  │
  ▼
fact tables
  │
  ▼
Reports engine
  │
  ├──▶ player reports
  └──▶ administrator reports
```

See the [glossary](glossary.md) for the definition of an actor.
