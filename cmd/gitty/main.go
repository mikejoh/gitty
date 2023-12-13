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
	"github.com/mikejoh/gity/internal/buildinfo"
)

type gittyOptions struct {
	path    string
	version bool
}

func main() {
	var gittyOpts gittyOptions
	flag.StringVar(&gittyOpts.path, "path", ".", "The path to scan for branches.")
	flag.BoolVar(&gittyOpts.version, "version", false, "Print the version number.")
	flag.Parse()

	if gittyOpts.version {
		fmt.Println(buildinfo.Get())
		os.Exit(0)
	}

	var gitRepoPaths []string
	currentTime := time.Now()

	err := filepath.Walk(gittyOpts.path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			gitRepoPaths = append(gitRepoPaths, filepath.Dir(path))
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

		var lastCommitSinceNow string

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

		if path == "." {
			path, err = os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
		}

		rows = append(rows, table.Row{
			path,
			ref.Name().Short(),
			commit.Author.When,
			lastCommitSinceNow,
			joinedRemotes,
		})
	}

	tw := table.NewWriter()
	header := table.Row{"directory", "branch", "last commit", "last commit age", "remotes"}
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

// // Fetch the latest updates from origin
// err = repo.Fetch(&git.FetchOptions{})
// if errors.Is(err, git.NoErrAlreadyUpToDate) {
// 	log.Printf("%s is already up to date!", path)
// } else if err != nil {
// 	log.Fatal(err)
// }

// // Get the commit history from HEAD
// cIter, err := repo.Log(&git.LogOptions{From: head.Hash()})
// if err != nil {
// 	log.Fatal(err)
// }

// localCount, originCount := 0, 0

// // Count local commits
// err = cIter.ForEach(func(c *object.Commit) error {
// 	localCount++
// 	return nil
// })
// if err != nil {
// 	log.Fatal(err)
// }

// // Get the reference for origin/HEAD
// originRef, err := repo.Reference("refs/remotes/origin/main", true)
// if err != nil {
// 	// If origin/main does not exist, try origin/master
// 	originRef, err = repo.Reference("refs/remotes/origin/master", true)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// // Get the commit history from origin/HEAD
// cIter, err = repo.Log(&git.LogOptions{From: originRef.Hash()})
// if err != nil {
// 	log.Fatal(err)
// }
// // Count origin commits
// err = cIter.ForEach(func(c *object.Commit) error {
// 	originCount++
// 	return nil
// })
// if err != nil {
// 	log.Fatal(err)
// }

// // Get the remote HEAD reference
// remoteRef, err := r.Reference(plumbing.NewRemoteReferenceName("origin", "master"), true)
// if err != nil {
// 	fmt.Println(err)
// 	return
// }

// // Create the commit iterators
// localIter, err := r.Log(&git.LogOptions{From: headRef.Hash()})
// if err != nil {
// 	fmt.Println(err)
// 	return
// }

// remoteIter, err := r.Log(&git.LogOptions{From: remoteRef.Hash()})
// if err != nil {
// 	fmt.Println(err)
// 	return
// }

// // Count the number of commits
// localCount, remoteCount := 0, 0
// localIter.ForEach(func(c *object.Commit) error {
// 	localCount++
// 	return nil
// })
// remoteIter.ForEach(func(c *object.Commit) error {
// 	remoteCount++
// 	return nil
// })

// fmt.Printf("Number of local commits: %d\n", localCount)
// fmt.Printf("Number of remote commits: %d\n", remoteCount)
