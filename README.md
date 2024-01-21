# gitty

<p align="center">
<img src="https://github.com/mikejoh/gitty/assets/899665/57b83aee-3f75-4cde-a6ad-e9e052f0d6ba" alt="gitty" />
</p>

`gitty` - Output some nice-to-have information about your local git repositories.

## Install

1. `make build`
2. `make install`

## Usage

Point at a specific directory:
```
gitty -path ./trela/
┌───────────┬────────┬───────────────────────────┬─────────────────┐
│ DIRECTORY │ BRANCH │ LAST COMMIT               │ LAST COMMIT AGE │
├───────────┼────────┼───────────────────────────┼─────────────────┤
│ trela     │ main   │ 2023-12-13T14:42:42+01:00 │ 0d 0h 5m        │
└───────────┴────────┴───────────────────────────┴─────────────────┘
```
Iterate through branches in dir '.':
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
