# gitty

<p align="center">
<img src="https://github.com/mikejoh/gitty/assets/899665/57b83aee-3f75-4cde-a6ad-e9e052f0d6ba" alt="gitty" />
</p>

`gitty` - A CLI tool to display useful information about your local Git repositories.

---

## Features

- Display branch information, including the last commit, age, author, and remotes.
- Filter commits by author.
- Traverse all branches in a repository.
- Skip specific directories during scanning.
- Truncate long branch names for better readability.

---

## Install

### From source

1. `git clone https://github.com/mikejoh/gitty.git`
2. `cd gitty`
3. `make build`
4. `make install` (assumes thath `~/.local/bin` is used)

### Download and run

1. Download (using `v0.1.3` as an example):

```bash
curl -LO https://github.com/mikejoh/gitty/releases/download/0.1.3/gitty_0.1.3_linux_amd64.tar.gz
```

2. Unpack:

```bash
tar xzvf gitty_0.1.3_linux_amd64.tar.gz
```

3. Run:

```bash
./gitty --version
```

## Usage

```bash
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

```bash
gitty --path .
┌────────────┬────────┬────────┬─────────────────────────────────┬──────────┬──────────────────┬─────────┐
│ REPOSITORY │ SIZE   │ BRANCH │ LAST COMMIT                     │ AGE      │ AUTHOR           │ REMOTES │
├────────────┼────────┼────────┼─────────────────────────────────┼──────────┼──────────────────┼─────────┤
│ gitty      │ 6.67MB │ main   │ 2025-12-04 11:23:21 +0100 +0100 │ 0d 0h 1m │ Mikael Johansson │ origin  │
└────────────┴────────┴────────┴─────────────────────────────────┴──────────┴──────────────────┴─────────┘
```

Iterate through all branches in local repository:

```bash
gitty --path . --branches
┌────────────┬──────┬─────────────────────────┬─────────────────────────────────┬──────────────┬──────────────────┬─────────┐
│ REPOSITORY │ SIZE │ BRANCH                  │ LAST COMMIT                     │ AGE          │ AUTHOR           │ REMOTES │
├────────────┼──────┼─────────────────────────┼─────────────────────────────────┼──────────────┼──────────────────┼─────────┤
│ gitty      │ N/A  │ calc-disk-size-per-path │ 2025-01-13 08:49:02 +0100 +0100 │ 325d 2h 36m  │ Mikael Johansson │ origin  │
│ gitty      │ N/A  │ iterate-on-branches     │ 2024-01-21 21:08:03 +0100 +0100 │ 682d 14h 17m │ Mikael Johansson │ origin  │
│ gitty      │ N/A  │ main                    │ 2025-12-04 11:23:21 +0100 +0100 │ 0d 0h 1m     │ Mikael Johansson │ origin  │
└────────────┴──────┴─────────────────────────┴─────────────────────────────────┴──────────────┴──────────────────┴─────────┘
```
