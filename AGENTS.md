# Alpha Workflow

- Work directly on `main` while the project version is a pre-release with the
  `alpha` identifier. Do not create a branch unless the user asks for one.
- Follow semantic versioning for every change.
- Committing to `main` and pushing are pre-approved while the project is in
  alpha, provided the change follows semantic versioning.
- Whenever a change bumps the project version, update `version.go` in the same
  commit.
- Run the relevant tests and formatting checks before committing and pushing.

# Database Safety

- Database create and initialize commands must never create missing directory
  paths. Treat any missing path component as a hard error.
- Database commands that open an existing persistent store must fail if any
  component of its directory path is missing.
- The open command must verify the SQLite application ID and refuse to open a
  database created by another application.
- Use the ASCII bytes for `TNYC` (`0x544E5943`) as the SQLite application ID.
- Use ZombieZen's SQLite driver and `sqlitemigration` package.
- Enable write-ahead logging (WAL) and foreign-key enforcement on every
  connection to a persistent store. Enable foreign-key enforcement on every
  connection to a temporary store.
- Verify the migration version with the ZombieZen `sqlitemigration` package's
  PRAGMA check. Do not implement a custom migration-version check.
- In-memory, non-persistent stores must support both a single shared instance
  and a private instance per request so tests can use either mode.

# Authentication

- Use the bcrypt package for password authentication. Use `bcrypt.MinCost`
  while the project version is an alpha pre-release.
