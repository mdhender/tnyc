# Entity reference

An **Entity** is an individually addressable, indivisible object in one game.
Orders generally name one Entity as their subject or target. Each Entity has a
unique number and name within its game.

This model is based primarily on sections 5.1–5.3 of the
[T'Nyc rules](rules.md), with additional attributes needed to persist and
process a strategic turn-based game. Attributes marked **rules** are stated or
directly implied by the rules. Attributes marked **model** extend the original
rules for implementation.

## Kinds

Every Entity has exactly one `kind`.

| Kind | Meaning | Examples from the rules |
| --- | --- | --- |
| `person` | One independently acting character. | Player character, follower, non-player character. |
| `group` | Multiple similar beings that act as one indivisible unit. | Troop, mercenary company, undead unit. |
| `place` | A geographic or constructed location that can be governed. | Province, city, castle, fortress. |
| `vehicle` | A mobile conveyance that can carry a stack. | Ship. |
| `thing` | An item or fungible lot that cannot act independently. | Artifact, load of trade goods, horses. |
| `event` | A persistent game-world condition with its own identity. | Curse, quest. |

`person` and `group` refine the rules' broad use of **unit**. Rules text also
occasionally calls things “units” for order-writing purposes; that usage does
not change a thing's Entity kind.

An Entity's kind determines which attributes and orders are valid. Kind is
stable for the Entity's lifetime. If an effect transforms an Entity into a
fundamentally different kind, the implementation records the transformation
rather than silently reinterpreting its existing attributes.

## Core attributes

Every Entity has the identity and lifecycle attributes below. The other
attributes are optional relationships or values whose applicability depends on
kind and current state.

| Attribute | Source | Meaning |
| --- | --- | --- |
| `id` | rules | Immutable Entity number, unique within the game. |
| `game-id` | model | Game containing the Entity. |
| `name` | rules | Entity name, unique within the game. |
| `kind` | model | One of the kinds listed above. |
| `status` | model | Lifecycle state, such as `active`, `inactive`, `destroyed`, or `completed`. |
| `location-id` | rules | Place containing the Entity. A top-level province has no containing place. |
| `stack-id` | rules | Stack position occupied by the Entity, if any. |
| `controlled-by-id` | model | Actor currently authorized to submit orders for the Entity, if it can receive orders. |
| `lord-entity-id` | rules | Entity to which this Entity has sworn allegiance; absent for things and events. |
| `loyalty` | rules | Degree of loyalty to the lord, when allegiance applies. |
| `default-attitude` | rules | Default reaction to other Entities: `friendly`, `neutral`, or `hostile`. |
| `opinions` | rules | Per-Entity reactions that override the default attitude. |
| `treasury` | rules | Silver held by the Entity. |
| `skills` | rules | Skill ratings possessed by the Entity. Ratings ordinarily range from 1 to 10. |
| `effects` | model | Temporary or persistent modifiers currently affecting the Entity. |
| `created-turn` | model | Turn in which the Entity entered play. |
| `ended-turn` | model | Turn in which it left play, if applicable. |

`controlled-by-id` and `lord-entity-id` are different relationships. Control
answers which player or engine agent may issue orders. Lordship describes the
in-world chain of allegiance. A player always controls their primary character,
but persuasion, terror, swearing, captivity, and similar mechanics can change
other relationships.

The rules limit an Entity to five direct followers. The transitive lordship
tree is its faction. An implementation must prevent lordship cycles and enforce
the direct-follower limit for kinds that can control followers.

## Person and group attributes

| Attribute | Source | Meaning |
| --- | --- | --- |
| `species-id` | rules | Species shared by members of the Entity. |
| `headcount` | rules | Number of members; always 1 for a person. |
| `armor` | rules | Armor rating from 1 through 9. |
| `health` | model | Current capacity to survive wounds and continue acting. |
| `captor-entity-id` | rules | Entity holding this Entity prisoner, if captured. |
| `maintenance` | rules | Periodic support cost. |
| `salary` | rules | Contractual pay, principally for mercenary groups. |
| `movement-mode` | rules | Means of movement supplied by traits or skills, such as foot, mounted, swimming, or flying. |

Skills belong to the Entity, not to each member of a group. Recruiting into a
group adds members of the same species and skill profile. Armor and movement
affect travel and combat, while species and the mix of combat skills can alter
combat strength by terrain.

## Place attributes

| Attribute | Source | Meaning |
| --- | --- | --- |
| `place-kind` | rules | Province, city, castle, fortress, or another location category. |
| `terrain-id` | rules | Terrain classification for the place. |
| `population` | rules | Resident population, expressed as families where the rules do so. |
| `fortification` | rules | Defensive wall or fortification rating. |
| `routes` | rules | Directed exits to other places, including direction, route terrain, and travel time. |
| `production` | rules | Commodities produced, including seasonal availability and observed price data. |
| `consumption` | rules | Commodities demanded by the local market. |
| `economy` | model | Economic condition used for work, trade, production, and tax calculations. |
| `tax-base` | model | Taxable capacity derived from population and economy. |
| `damage` | rules | Degradation caused by pillaging, terror, combat, or similar effects. |
| `garrisons` | rules | Entities assigned to protect the place. |

A place may contain other places. For example, a city or fortress can be
entered from its surrounding province. Routes describe strategic-map
adjacency; containment describes where an Entity currently is.

## Vehicle attributes

| Attribute | Source | Meaning |
| --- | --- | --- |
| `vehicle-kind` | model | Vehicle category, such as ship. |
| `capacity` | model | Limit on carried people, groups, and things. |
| `condition` | model | Current structural state. |
| `movement-mode` | rules | Terrain or medium the vehicle can traverse. |
| `captain-entity-id` | rules | Entity whose movement order and relevant skills direct the vehicle. |

The rules describe a ship's captain as the Entity at the top of its stack.
`captain-entity-id` is therefore derived from stack order unless later rules
introduce a separate assignment.

## Thing attributes

| Attribute | Source | Meaning |
| --- | --- | --- |
| `thing-kind` | model | Commodity, artifact, equipment, mount, or another item category. |
| `quantity` | rules | Number of fungible items represented by the Entity. |
| `unit-weight` | model | Encumbrance or cargo use per item. |
| `condition` | model | Current physical state, when durability matters. |
| `powers` | rules | Special capabilities of an artifact or magical item. |
| `commodity-id` | rules | Traded commodity represented by the Entity, when applicable. |

Things cannot move by themselves or control followers. They must be stacked
under a person to move or be sold. Only orders explicitly valid for things,
such as stack-management orders, may be issued in their name.

## Event attributes

| Attribute | Source | Meaning |
| --- | --- | --- |
| `event-kind` | rules | Persistent condition category, such as curse or quest. |
| `source-entity-id` | model | Entity that created the event, if known. |
| `target-entity-id` | model | Entity directly affected by the event, if any. |
| `starts-turn` | model | First turn in which the event is effective. |
| `expires-turn` | model | Turn after which the event ends, if finite. |
| `state` | model | Progress or resolution state. |
| `conditions` | model | Requirements for advancement, removal, success, or failure. |

Events are addressable because the rules identify certain persistent events as
Entities. An event may affect a place, faction, or stack through relationships
without becoming that object.

## Observation and hidden information

Entity attributes describe authoritative game state. What an actor knows about
an Entity is a separate, turn-stamped observation. Reports can omit an Entity,
hide its lord, or expose uncertain skills, attitudes, market data, and other
attributes according to scouting, stealth, lore, and similar mechanics. An
observation must not overwrite the underlying Entity with incomplete or stale
information.

## Entity and stack distinction

A [stack](glossary.md#stack) is an ordered grouping, not an Entity. It has no
Entity number, name, loyalty, or independent orders. The Entity at the top is
the stack leader and can represent the stack when issuing movement and combat
orders. Persuasion and other single-target effects still apply to one Entity,
not to the stack as a whole.

Places, vehicles, and things remain Entities when they contain or carry a
stack. Do not represent a ship, settlement, or other game-world object as a
stack kind.
