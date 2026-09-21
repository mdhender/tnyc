# tnyc
An unfaithful remake of a lost(?) classic.

## Documentation

- [Data model](docs/data-model.md)
- [Engines](docs/engines.md)
- [Entity reference](docs/references/entity.md)
- [Glossary](docs/references/glossary.md)
- [T'Nyc rules](docs/references/rules.md)

## Importing a world

Validate the WGVC schema-v1 export and save the T'Nyc map as JSON:

```sh
go run ./cmd/tnyc world
```

The command reads `var/wgvc-export.json` and writes `var/tnyc-world.json` by
default. Use `--input` and `--output` to select other paths.

## Creating a database

Create a datastore in an existing writable directory using an imported world:

```sh
go run ./cmd/tnyc database create
```

The command defaults to `--db-path var/` and
`--world-map var/tnyc-world.JSON`. It creates `tnyc.json` in the database path
and refuses to overwrite an existing database.

## Acknowledgements

The [T'Nyc rules](docs/references/rules.md) are copied from the
[Olympia PBM Archive](https://www.pbm.com/oly/archive/tnyc/rules). To the best
of our understanding, Rich Skrenta released the rules into the public domain.
