package commands

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func (r *GitRepoInfo) SetRepoPath() (err error) {
	// set current path in order to execute git commands
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	r.Path = path
	return nil
}

func (r *GitRepoInfo) IsAGitRepo() (err error) {
	// check if it is a git repo
	// git rev-parse --is-inside-work-tree
	cmd := exec.Command("git", "-C", r.Path, "rev-parse", "--is-inside-work-tree")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("%s %w", stderr.String(), err)
		return err
	}
	return nil
}

func (r *GitRepoInfo) HasRemote() (err error) {
	// check if repo has remote origin
	// git remote -v
	r.HasOrigin = false
	cmd := exec.Command("git", "-C", r.Path, "remote", "-v")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("%w %s", err, stderr.String())
		return err
	}
	if strings.Contains(stdout.String(), "origin") {
		r.HasOrigin = true
		return nil
	}

	err = fmt.Errorf("no origin found")
	return err
}

func (r *GitRepoInfo) HasTags() (err error) {
	// check if repo has tags
	// git describe --tags --abbrev=0
	r.HasTag = false
	cmd := exec.Command("git", "-C", r.Path, "describe", "--tags", "--abbrev=0")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("no tags found!\n[*] You should create a tag first with 'taggo init'")

		return err
	}
	r.HasTag = true
	return nil
}

func (r *GitRepoInfo) CurrentTag() (tag string, err error) {
	// get current tag
	// git tag --sort=committerdate ==> linux, macos
	// git describe --tags --abbrev=0 ==> windows
	err = r.HasTags()
	if err != nil {
		return "", err
	}

	// FIXME: this is a workaround for windows
	// more tags on a commit will return the first one, not the last one
	var stdout, stderr bytes.Buffer
	if runtime.GOOS == "windows" {
		cmd := exec.Command("git", "-C", r.Path, "describe", "--tags", "--abbrev=0")
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			err = fmt.Errorf("%w %s", err, stderr.String())
			return "", err
		}
	} else {
		cmd1 := exec.Command("git", "-C", r.Path, "tag", "--sort=committerdate")
		cmd2 := exec.Command("tail", "-1")
		cmd2.Stdin, _ = cmd1.StdoutPipe()
		cmd2.Stdout = &stdout
		cmd2.Stderr = &stderr
		err1 := cmd1.Start()
		if err1 != nil {
			fmt.Println(err1)
			os.Exit(1)
		}
		err2 := cmd2.Run()
		if err2 != nil {
			fmt.Println(err1)
			os.Exit(1)
		}
	}

	tag = strings.TrimSpace(stdout.String())

	r.LastTag = tag
	return tag, nil
}

func (r *GitRepoInfo) Prerequisites() (err error) {

	err = r.SetRepoPath()
	if err != nil {
		return err
	}

	err = r.IsAGitRepo()
	if err != nil {
		return err
	}

	err = r.HasRemote()
	if err != nil {
		return err
	}

	// err = r.HasTags()
	// if err != nil {
	// 	return err
	// }

	return nil

}

func (r *GitRepoInfo) CreateTag(newTag string) (err error) {
	// tag current commit
	// git tag <tag>
	cmd := exec.Command("git", "-C", r.Path, "tag", newTag)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("%w %s", err, stderr.String())
		return err
	}
	return nil
}

func (r *GitRepoInfo) PushTag(newTag string) (err error) {
	// push tag to remote
	// git push origin <tag>
	cmd := exec.Command("git", "-C", r.Path, "push", "origin", newTag)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("%w %s", err, stderr.String())
		return err
	}
	return nil
}

func (r *GitRepoInfo) DeleteLastTagOnLocal() (err error) {
	// delete tag locally
	// git tag -d <tag>
	cmd := exec.Command("git", "-C", r.Path, "tag", "-d", r.LastTag)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("%w %s", err, stderr.String())
		return err
	}
	return nil
}

func (r *GitRepoInfo) DeleteLastTagOnRemote() (err error) {
	// delete tag remotely
	// git push origin --delete <tag>
	// if r.LastTag == "v0.0.0" {
	// 	// tag v0.0.0 not recognized by github
	// 	return nil
	// }
	cmd := exec.Command("git", "-C", r.Path, "push", "origin", "--delete", r.LastTag)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("%w %s", err, stderr.String())
		return err
	}
	return nil
}

func (r *GitRepoInfo) InitRepo() (err error) {
	// create tag v0.0.0
	err = r.Prerequisites()
	if err != nil {
		return err
	}

	r.HasTags() // ignore error

	if !r.HasTag {
		fmt.Println("[*] Initializing git repo")
		firstTag := "v0.0.0"
		err = r.CreateTag(firstTag)
		if err != nil {
			return err
		}
		err = r.PushTag(firstTag)
		if err != nil {
			return err
		}
		fmt.Println("[*] Added tag v0.0.0")
		return nil
	} else {
		fmt.Println("[*] Repo already initialized")
	}
	return nil
}
