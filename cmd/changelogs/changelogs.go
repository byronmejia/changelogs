package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"regexp"
	"strings"
)

func main() {
	repository, err := git.PlainOpen("./")
	if err != nil {
		panic(err)
	}

	tagRefs, err := repository.Tags()

	if err != nil {
		panic(err)
	}

	var mostRecentTagCommit *object.Commit

	err = tagRefs.ForEach(func(reference *plumbing.Reference) error {
		tagObject, err := repository.TagObject(reference.Hash())
		if err != nil {
			return err
		}

		commit, err := tagObject.Commit()
		if err != nil {
			return err
		}

		if mostRecentTagCommit == nil {
			mostRecentTagCommit = commit
			return nil
		}

		if commit.Author.When.Sub(mostRecentTagCommit.Author.When) >= 0 {
			mostRecentTagCommit = commit
		}

		return nil
	})
	if err != nil {
		panic(err)
	}

	headRef, err := repository.Head()
	if err != nil {
		panic(err)
	}

	headCommit, err := repository.CommitObject(headRef.Hash())
	if err != nil {
		panic(err)
	}

	nextCommit := headCommit

	fmt.Println(nextCommit.NumParents())

	log, err := repository.Log(&git.LogOptions{
		From: headRef.Hash(),
		Order: git.LogOrderCommitterTime,
	})

	if err != nil {
		panic(err)
	}

	allCommitTitles := make([]string, 1)

	for nextCommit.ID() != mostRecentTagCommit.ID() {
		currentCommit, err := log.Next()
		if err != nil {
			panic(err)
		}
		nextCommit = currentCommit
		commitMessageTitle := strings.Split(currentCommit.Message, "\n")[0]
		allCommitTitles = append(allCommitTitles, commitMessageTitle)
	}

	angularPattern, err := regexp.Compile("^(\\w*)(?:\\((.*)\\))?: (.*)$")

	if err != nil {
		panic(err)
	}

	for _, v := range allCommitTitles {
		fmt.Println(v)
		fmt.Println(angularPattern.MatchString(v))
	}
}
