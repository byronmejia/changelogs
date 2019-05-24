package main

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"regexp"
	"strings"
)

type commitCorrespondence struct {
	scope string
	subject string
}

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

	log, err := repository.Log(&git.LogOptions{
		From: headRef.Hash(),
		Order: git.LogOrderCommitterTime,
	})

	if err != nil {
		panic(err)
	}

	allCommitTitles := make([]string, 0)

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

	goodCommitTitles := make([]string, 0)
	junkCommitTitles := make([]string, 0)

	goodCommitMap := make(map[string][]commitCorrespondence)

	for _, v := range allCommitTitles {
		if angularPattern.MatchString(v) {
			goodCommitTitles = append(goodCommitTitles, v)
			for _, v := range angularPattern.FindAllStringSubmatch(v, 4) {
				if goodCommitMap[strings.ToLower(v[1])] == nil {
					goodCommitMap[strings.ToLower(v[1])] = make([]commitCorrespondence, 0)
				}
				goodCommitMap[strings.ToLower(v[1])] = append(goodCommitMap[strings.ToLower(v[1])], commitCorrespondence{
					scope: v[2],
					subject: v[3],
				})
			}
		} else {
			junkCommitTitles = append(junkCommitTitles, v)
		}
	}

	for k, v := range goodCommitMap {
		fmt.Println(k)
		scopes := make(map[string][]string)
		for _, v := range v {
			if scopes[v.scope] == nil {
				scopes[v.scope] = make([]string, 0)
			}

			scopes[v.scope] = append(scopes[v.scope], v.subject)
		}
		for k, v := range scopes {
			fmt.Print("\t")
			fmt.Println(k)
			for _, v := range v {
				fmt.Print("\t\t")
				fmt.Println(v)
			}
		}
	}
}
