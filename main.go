package main

import (
	"fmt"
	"os"

	"gopkg.in/src-d/go-git.v4"
	. "gopkg.in/src-d/go-git.v4/_examples"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func checkIfError(args ...interface{}) {
	if len(args) < 2 || args[1] == nil {
		return
	}
	fmt.Println(args...)
	os.Exit(1)
}

// Example of how to:
// - Clone a repository into memory
// - Get the HEAD reference
// - Using the HEAD reference, obtain the commit this reference is pointing to
// - Using the commit, obtain its history and print it
func main() {
	// Clones the given repository, creating the remote, the local branches
	// and fetching the objects, everything in memory:
	Info("git clone https://github.com/dmigwi/golang-modules.git")
	r, err := git.PlainClone("data", false, &git.CloneOptions{
		URL: "https://github.com/dmigwi/golang-modules.git",
	})

	if err == git.ErrRepositoryAlreadyExists {
		r, err = git.PlainOpen("data")
		checkIfError(err)

		w, err := r.Worktree()
		checkIfError(err)

		// Pull the latest changes from the origin remote and merge into the current branch
		Info("git pull origin")
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	}

	checkIfError(err)

	// Gets the HEAD history from HEAD, just like this command:
	Info("git log")

	// ... retrieves the branch pointed by HEAD
	ref, err := r.Head()
	checkIfError("r.Head()...", err)

	// ... retrieves the commit history
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	// cIter, err := r.Object(plumbing.CommitObject, nil)
	checkIfError("r.Log(...", err)

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		fmt.Println(c)

		fromTree, err := c.Tree()
		if err != nil {
			return err
		}

		toTree := &object.Tree{}
		if c.NumParents() != 0 {
			firstParent, err := c.Parents().Next()
			if err != nil {
				return err
			}

			toTree, err = firstParent.Tree()
			if err != nil {
				return err
			}
		}

		patch, err := toTree.Patch(fromTree)
		if err != nil {
			return err
		}
		fmt.Printf("%+v", patch)
		fmt.Println("\n\n")
		return nil
	})

	checkIfError("cIter.ForEach...", err)
}
