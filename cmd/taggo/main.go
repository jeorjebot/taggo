package main

import (
	"fmt"
	"os"

	"github.com/jeorjebot/taggo/pkg/commands"
	flags "github.com/jessevdk/go-flags"
)

var (
	opts   Options
	repo   commands.GitRepoInfo = commands.GitRepoInfo{}
	parser *flags.Parser        = flags.NewParser(&opts, flags.Default)
)

// Options holds the CLI args
type Options struct {
	Tag   string `short:"t" long:"tag" description:"Tag to create"`
	Major bool   `short:"M" long:"major" description:"Bump major version"`
	Minor bool   `short:"m" long:"minor" description:"Bump minor version"`
	Patch bool   `short:"p" long:"patch" description:"Bump patch version"`
	//PreReleaseName string `short:"n" long:"pre-release-name" description:"create a pre-release tag"`
	Delete bool `short:"d" long:"delete" description:"Delete last tag"`
}

func init() {

	_, err := parser.Parse()
	if err != nil {
		if err.(*flags.Error).Type == flags.ErrHelp {
			os.Exit(0)
		}
		fmt.Println("==> Error parsing flags: " + err.Error())
		os.Exit(1)
	}
}

func main() {
	// if there are no tags, init repo
	if len(os.Args) > 1 {
		if os.Args[1] == "init" {
			if len(os.Args) > 2 {
				fmt.Println("[*] Error: no arguments expected for command 'init'")
				fmt.Println("[*] Usage: taggo init")
				os.Exit(1)
			} else {
				err := repo.InitRepo()
				checkError(err)
				os.Exit(0)
			}
		}
	}

	err := repo.Prerequisites()
	checkError(err)

	if opts.Tag != "" {
		// check tag format vX.Y.Z
		err = repo.CheckTagFormat(opts.Tag)
		checkError(err)

		err = repo.CreateTag(opts.Tag)
		checkError(err)

		fmt.Println("[*] New tag: " + opts.Tag)

		err = repo.PushTag(opts.Tag)
		checkError(err)

		fmt.Println("[*] Tag pushed successfully")
		os.Exit(0)
	}

	tag, err := repo.CurrentTag()
	checkError(err)

	fmt.Println("[*] Current tag: " + tag)

	// if no options, show current tag and exit
	if len(os.Args) == 1 {
		os.Exit(0)
	}

	// delete last tag if needed
	if opts.Delete {
		fmt.Println("[*] Deleting tag " + tag)
		err = repo.DeleteLastTagOnLocal()
		checkError(err)

		err = repo.DeleteLastTagOnRemote()
		checkError(err)

		fmt.Println("[*] Tag deleted successfully")
		os.Exit(0)
	}

	// increment tag based on options
	newTag := ""
	if opts.Major {
		newTag, err = repo.IncMajor()
		checkError(err)
	}
	if opts.Minor {
		newTag, err = repo.IncMinor()
		checkError(err)
	}
	if opts.Patch {
		newTag, err = repo.IncPatch()
		checkError(err)
	}

	if newTag != "" {
		fmt.Println("[*] New tag: " + newTag)
	}

	// create new tag
	err = repo.CreateTag(newTag)
	checkError(err)

	// push new tag
	err = repo.PushTag(newTag)
	checkError(err)

	fmt.Println("[*] Tag pushed successfully")

}

func checkError(err error) {
	if err != nil {
		fmt.Println("==> Error: " + err.Error())
		os.Exit(1)
	}
}
