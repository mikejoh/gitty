package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mikejoh/gitty/internal/buildinfo"
)

const maxBranchNameLength = 25

type gittyOptions struct {
	path               string
	branches           bool
	version            bool
	skipDirs           string
	truncateBranchName bool
}

var skipDirs = []string{
	".terraform",
}

func main() {
	var gittyOpts gittyOptions
	flag.StringVar(&gittyOpts.path, "path", ".", "The path to scan for branches.")
	flag.StringVar(&gittyOpts.skipDirs, "skip-dirs", "", "Comma separated list of directories to skip.")
	flag.BoolVar(&gittyOpts.branches, "branches", false, "Loop through all branches and print the last commit age for each branch.")
	flag.BoolVar(&gittyOpts.truncateBranchName, "truncate", false, "Truncate branch names to 25 characters.")
	flag.BoolVar(&gittyOpts.version, "version", false, "Print the version number.")
	flag.Parse()

	if gittyOpts.version {
		log.Println(buildinfo.Get())
		os.Exit(0)
	}

	if gittyOpts.skipDirs != "" {
		sd := strings.Split(gittyOpts.skipDirs, ",")
		skipDirs = append(skipDirs, sd...)
	}

	var gitRepoPaths []string
	currentTime := time.Now()

	err := filepath.Walk(gittyOpts.path, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			for _, dir := range skipDirs {
				if fileInfo.Name() == dir {
					return filepath.SkipDir
				}
			}
			if fileInfo.Name() == ".git" {
				gitRepoPaths = append(gitRepoPaths, filepath.Dir(path))
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	rows := []table.Row{}

	for _, path := range gitRepoPaths {
		repo, err := git.PlainOpen(path)
		if errors.Is(err, git.ErrRepositoryNotExists) {
			log.Fatal(err)
		}

		remotes, err := repo.Remotes()
		if err != nil {
			log.Fatal(err)
		}

		var remoteNames []string
		for _, remote := range remotes {
			remoteNames = append(remoteNames, remote.Config().Name)
		}

		joinedRemotes := strings.Join(remoteNames, ", ")

		var lastCommitSinceNow string

		if gittyOpts.branches {
			branchIter, err := repo.Branches()
			if err != nil {
				log.Fatal(err)
			}

			err = branchIter.ForEach(func(branch *plumbing.Reference) error {
				// Skip if no commits exists
				_, err := repo.CommitObject(branch.Hash())
				if err != nil {
					return nil
				}

				// Skip detached HEAD
				if branch.Name() == plumbing.HEAD {
					return nil
				}

				var lastCommitSinceNow string

				commit, err := repo.CommitObject(branch.Hash())
				if err == nil {
					duration := currentTime.Sub(commit.Author.When)
					days := duration.Hours() / 24
					hours := int(duration.Hours()) % 24
					minutes := int(duration.Minutes()) % 60
					lastCommitSinceNow = fmt.Sprintf("%dd %dh %dm", int(days), hours, minutes)
				} else if err != nil {
					log.Fatal(err)
				}

				if gittyOpts.truncateBranchName {
					branchName := branch.Name().Short()
					if len(branchName) > maxBranchNameLength {
						branchName = fmt.Sprintf("%.20s...", branchName)
					}

					rows = append(rows, table.Row{
						branchName,
						commit.Author.When,
						lastCommitSinceNow,
						commit.Author.Name,
						joinedRemotes,
					})
					return nil
				}

				rows = append(rows, table.Row{
					branch.Name().Short(),
					commit.Author.When,
					lastCommitSinceNow,
					commit.Author.Name,
					joinedRemotes,
				})
				return nil
			})
			if err != nil {
				log.Fatal(err)
			}

			continue
		}

		head, err := repo.Head()
		if err != nil {
			continue
		}

		// Skip if no commits exists
		_, err = repo.CommitObject(head.Hash())
		if err != nil {
			continue
		}

		// Skip detached HEAD
		if head.Name() == plumbing.HEAD {
			continue
		}

		ref, err := repo.Reference(plumbing.HEAD, true)
		if err != nil {
			log.Fatalf("%s: %s", path, err)
		}

		commit, err := repo.CommitObject(ref.Hash())
		if err == nil {
			duration := currentTime.Sub(commit.Author.When)
			days := duration.Hours() / 24
			hours := int(duration.Hours()) % 24
			minutes := int(duration.Minutes()) % 60
			lastCommitSinceNow = fmt.Sprintf("%dd %dh %dm", int(days), hours, minutes)
		} else if err != nil {
			log.Fatal(err)
		}

		rows = append(rows, table.Row{
			ref.Name().Short(),
			commit.Author.When,
			lastCommitSinceNow,
			commit.Author.Name,
			joinedRemotes,
		})

	}

	tw := table.NewWriter()
	header := table.Row{
		"branch",
		"last commit",
		"age",
		"author",
		"remotes",
	}
	tw.AppendHeader(header)
	tw.AppendRows(rows)
	tw.Style().Title.Align = text.AlignCenter
	tw.SetStyle(table.StyleLight)
	tw.SetOutputMirror(os.Stdout)
	tw.SetColumnConfigs([]table.ColumnConfig{{
		Name:        "last commit",
		Transformer: text.NewTimeTransformer(time.RFC3339, nil),
	}})
	tw.SortBy([]table.SortBy{
		{Name: "last commit", Mode: table.Dsc},
	})
	tw.Render()
}
