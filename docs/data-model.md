# Data Model

This document defines the core records and integrity rules for accounts, games,
players, and game entities.

## Identifier policy

Every record's `id` is an integer that is unique within that record type. Once
assigned, an ID is never reused, including after its record is deleted or made
inactive. Foreign-key fields contain the ID of the referenced record and are
not themselves unique unless a constraint below says otherwise.

## Records

### Account

An account identifies a person who can access the system.

| Field | Meaning |
| --- | --- |
| `id` | Account ID. |
| `email` | Normalized email address. Must be unique. |
| `display-name` | Name shown to other users. |
| `is-active` | Whether the account may be used. |
| `is-admin` | Whether the account has administrator privileges. |

### Game

A game is one instance of play.

| Field | Meaning |
| --- | --- |
| `id` | Game ID. |
| `is-active` | Whether the game is active. |
| `current-turn` | The game's current turn. |

### Player

A player associates an account with a game. The pair (`game-id`, `account-id`)
must be unique, so an account can have at most one player record in a game.

| Field | Meaning |
| --- | --- |
| `id` | Player ID. |
| `game-id` | ID of the game. |
| `account-id` | ID of the participating account. |
| `is-active` | Whether the player is active in the game. |

### Entity

An entity is an individually addressable object in the game world.

| Field | Meaning |
| --- | --- |
| `id` | Entity ID. |
| `controlled-by-id` | ID of the actor that controls the entity. |
| `stack-id` | ID of the stack containing the entity. |

### Stack

A stack is an ordered group of entities and substacks at a location. Its kind
describes the group, such as a ship or settlement.

| Field | Meaning |
| --- | --- |
| `id` | Stack ID. |
| `location-id` | ID of the stack's location. |
| `kind` | Stack kind, such as `ship` or `settlement`. |
| `controlled-by-id` | ID of the actor that controls the stack. |

### Element

An element places either an entity or a substack at a particular position in a
parent stack.

| Field | Meaning |
| --- | --- |
| `stack-id` | ID of the parent stack. |
| `sequence` | Position within the parent stack. |
| `entity-id` | ID of the entity at that position, if the element contains an entity. |
| `sub-stack-id` | ID of the stack at that position, if the element contains a substack. |

Exactly one of `entity-id` and `sub-stack-id` must be present. Within a parent
stack, `sequence` identifies an element's position and must be unique. An entity
or substack can occupy at most one stack position. A stack cannot contain itself
directly or indirectly.

## Relationships

```text
Account 1 ─── 0..* Player 0..* ─── 1 Game
Actor   1 ─── 0..* Entity
Actor   1 ─── 0..* Stack
Stack   1 ─── 0..* Element 0..1 ─── 1 Entity
                         │
                         └── 0..1 ─── 1 Stack (substack)
```

An active player is an [actor](glossary.md#actor). Engines may also act through
agents.

## Open modeling decisions

The initial model does not yet define:

- the record and key used to represent actors and engine agents; or
- whether `Entity.stack-id` or an entity-valued `Element` is the authoritative
  source for entity membership in a stack.

These decisions must be resolved before implementing the corresponding foreign
keys and consistency constraints.
