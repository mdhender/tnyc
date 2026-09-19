# Alpha Workflow

- Work directly on `main` while the project version is a pre-release with the
  `alpha` identifier. Do not create a branch unless the user asks for one.
- Follow semantic versioning for every change.
- Committing to `main` and pushing are pre-approved while the project is in
  alpha, provided the change follows semantic versioning.
- Whenever a change bumps the project version, update `version.go` in the same
  commit.
- Run the relevant tests and formatting checks before committing and pushing.
