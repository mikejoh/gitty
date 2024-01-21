# gitty

`gitty` - Output some nice-to-have information about your local git repositories.

## Install

1. `make build`
2. `make install`

## Usage

Example:
```
gitty -path ./trela/
┌───────────┬────────┬───────────────────────────┬─────────────────┐
│ DIRECTORY │ BRANCH │ LAST COMMIT               │ LAST COMMIT AGE │
├───────────┼────────┼───────────────────────────┼─────────────────┤
│ trela     │ main   │ 2023-12-13T14:42:42+01:00 │ 0d 0h 5m        │
└───────────┴────────┴───────────────────────────┴─────────────────┘
```
```
gitty -branches
┌───────────┬────────┬───────────────────────────┬─────────────────┐
│ DIRECTORY │ BRANCH │ LAST COMMIT               │ LAST COMMIT AGE │
├───────────┼────────┼───────────────────────────┼─────────────────┤
│ trela     │ main   │ 2023-12-13T14:42:42+01:00 │ 0d 0h 5m        │
│ trela     │ test   │ 2023-12-18T12:32:42+01:00 │ 23d 3h 5m       │
└───────────┴────────┴───────────────────────────┴─────────────────┘
```

`gitty` defaults to find git repositories in the current working directory ('`.`').
