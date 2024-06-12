# gitty

<p align="center">
<img src="https://github.com/mikejoh/gitty/assets/899665/57b83aee-3f75-4cde-a6ad-e9e052f0d6ba" alt="gitty" />
</p>

`gitty` - Output some nice-to-have information about your local git repositories.

## Install

1. `make build`
2. `make install`

## Usage

```
Usage of gitty:
  -branches
    	Loop through all branches.
  -path string
    	The path to scan for branches. (default ".")
  -skip-dirs string
    	Comma separated list of directories to skip.
  -truncate
    	Truncate branch names to 25 characters.
  -version
    	Print the version number.
```

### Examples

Local git repository:
```
gitty
┌────────┬───────────────────────────┬────────────┬──────────────────┬─────────┐
│ BRANCH │ LAST COMMIT               │ AGE        │ AUTHOR           │ REMOTES │
├────────┼───────────────────────────┼────────────┼──────────────────┼─────────┤
│ main   │ 2024-05-22T13:10:25+02:00 │ 21d 7h 35m │ Mikael Johansson │ origin  │
└────────┴───────────────────────────┴────────────┴──────────────────┴─────────┘
```

Iterate through all branches in local repository:
```
gitty -branches
┌────────┬───────────────────────────┬────────────┬──────────────────┬─────────┐
│ BRANCH │ LAST COMMIT               │ AGE        │ AUTHOR           │ REMOTES │
├────────┼───────────────────────────┼────────────┼──────────────────┼─────────┤
│ main   │ 2024-05-22T13:10:25+02:00 │ 21d 7h 36m │ Mikael Johansson │ origin  │
│ test   │ 2024-05-22T13:10:25+02:00 │ 21d 7h 36m │ Mikael Johansson │ origin  │
└────────┴───────────────────────────┴────────────┴──────────────────┴─────────┘
```
