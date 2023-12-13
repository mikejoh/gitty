# gitty

`gitty` - Output some nice-to-have information about your local git repositories.

## Install

1. `make build`
2. `make install`

## Usage

Example:
```
gitty -path ./trela/
```
Output:
```
┌───────────┬────────┬───────────────────────────┬─────────────────┐
│ DIRECTORY │ BRANCH │ LAST COMMIT               │ LAST COMMIT AGE │
├───────────┼────────┼───────────────────────────┼─────────────────┤
│ trela     │ main   │ 2023-12-13T14:42:42+01:00 │ 0d 0h 5m        │
└───────────┴────────┴───────────────────────────┴─────────────────┘
```