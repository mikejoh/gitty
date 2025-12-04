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
	authorFilter       string // New flag for filtering by author
}

var skipDirs = []string{
	".terraform",
}

func main() {
	var gittyOpts gittyOptions
	flag.StringVar(&gittyOpts.path, "path", "", "The path to scan for branches (required).")
	flag.StringVar(&gittyOpts.skipDirs, "skip-dirs", "", "Comma separated list of directories to skip.")
	flag.BoolVar(&gittyOpts.branches, "branches", false, "List all branches.")
	flag.BoolVar(&gittyOpts.truncateBranchName, "truncate", false, "Truncate branch names to 25 characters.")
	flag.BoolVar(&gittyOpts.version, "version", false, "Print the version number.")
	flag.StringVar(&gittyOpts.authorFilter, "author", "", "Filter commits by author name.") // New flag
	flag.Parse()

	// Display usage if no arguments or flags are provided
	if flag.NFlag() == 0 && flag.NArg() == 0 {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if flag.NArg() > 0 {
		gittyOpts.path = flag.Arg(0)
	}

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
		size, err := getDirSize(path)
		if err != nil {
			log.Fatal(err)
		}

		repo, err := git.PlainOpen(path)
		if errors.Is(err, git.ErrRepositoryNotExists) {
			log.Fatal(err)
		}

		remotes, err := repo.Remotes()
		if err != nil {
			log.Fatal(err)
		}

		var originRemote string
		var remoteNames []string
		for _, remote := range remotes {
			if remote.Config().Name == "origin" {
				originRemote = remote.Config().URLs[0]
			}
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
				// Skip if no commits exist
				commit, err := repo.CommitObject(branch.Hash())
				if err != nil {
					return nil
				}

				// Skip detached HEAD
				if branch.Name() == plumbing.HEAD {
					return nil
				}

				// Filter by author
				if gittyOpts.authorFilter != "" && commit.Author.Name != gittyOpts.authorFilter {
					return nil
				}

				duration := currentTime.Sub(commit.Author.When)
				days := duration.Hours() / 24
				hours := int(duration.Hours()) % 24
				minutes := int(duration.Minutes()) % 60
				lastCommitSinceNow = fmt.Sprintf("%dd %dh %dm", int(days), hours, minutes)

				branchName := branch.Name().Short()
				if gittyOpts.truncateBranchName && len(branchName) > maxBranchNameLength {
					branchName = fmt.Sprintf("%.20s...", branchName)
				}

				rows = append(rows, table.Row{
					strings.TrimSuffix(filepath.Base(originRemote), ".git"),
					"N/A",
					branchName,
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

		// Skip if no commits exist
		commit, err := repo.CommitObject(head.Hash())
		if err != nil {
			continue
		}

		// Skip detached HEAD
		if head.Name() == plumbing.HEAD {
			continue
		}

		// Filter by author
		if gittyOpts.authorFilter != "" && commit.Author.Name != gittyOpts.authorFilter {
			continue
		}

		duration := currentTime.Sub(commit.Author.When)
		days := duration.Hours() / 24
		hours := int(duration.Hours()) % 24
		minutes := int(duration.Minutes()) % 60
		lastCommitSinceNow = fmt.Sprintf("%dd %dh %dm", int(days), hours, minutes)

		if originRemote == "" {
			originRemote = filepath.Base(path)
		}

		rows = append(rows, table.Row{
			strings.TrimSuffix(filepath.Base(originRemote), ".git"),
			fmt.Sprintf("%.2fMB", float64(size)/(1024*1024)),
			head.Name().Short(),
			commit.Author.When,
			lastCommitSinceNow,
			commit.Author.Name,
			joinedRemotes,
		})
	}

	tw := table.NewWriter()
	header := table.Row{
		"repository",
		"size",
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
	tw.Render()
}

func getDirSize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	return size, err
}
