# Glossary

## Actor

An **actor** submits orders. An actor is either:

- a **player**: an account participating actively in a game; or
- an **agent**: an engine component that creates and submits orders.

Actor is a role in order processing, not an [Entity](#entity). Its persistent
representation is not yet defined in the [data model](../data-model.md).

## Attitude

An **attitude** is an Entity's intended reaction to another Entity. The values
are **friendly**, **neutral**, and **hostile**. An Entity has a default attitude
and may declare a different opinion of a specific Entity.

## Entity

An **Entity** is a uniquely numbered and named, indivisible game-world object
that orders can address. Entity kinds and attributes are defined in the
[Entity reference](entity.md).

## Faction

A **faction** is the hierarchy formed by an Entity, its direct followers, and
their followers recursively. The Entity at the root is the ultimate lord.

## Lord

A **lord** is an Entity to which another Entity has sworn allegiance. An Entity
can have at most one direct lord and, under the original rules, no more than
five direct followers.

## Place

A **place** is an Entity that represents a geographic or constructed location,
such as a province, city, castle, or fortress.

## Stack

A **stack** is an ordered group of Entities that move and fight together. The
top Entity leads the stack. A stack is not itself an Entity and cannot be the
target of effects that apply to one Entity.

## Thing

A **thing** is an Entity such as an artifact or a load of trade goods. It
cannot move independently or control followers and must be carried in a stack.

## Unit

A **unit** is the rules' general term for an Entity used in order processing.
It most often means one person or a group of similar beings, but the original
rules sometimes also treat things as units. Use the precise Entity kind when
the distinction matters.
