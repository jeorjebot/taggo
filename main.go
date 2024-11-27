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
	Tag              string `short:"t" long:"tag" description:"Tag to create"`
	Major            bool   `short:"M" long:"major" description:"Bump major version"`
	Minor            bool   `short:"m" long:"minor" description:"Bump minor version"`
	Patch            bool   `short:"p" long:"patch" description:"Bump patch version"`
	PreReleaseName   string `short:"n" long:"pre-release" description:"create a pre-release tag"`
	Delete           bool   `short:"d" long:"delete" description:"Delete last tag"`
	Init             bool   `short:"i" long:"init" description:"Initialize git repo with taggo"`
	InitWithNoPrefix bool   `short:"I" long:"init-no-prefix" description:"Initialize git repo with taggo without 'v' prefix"`
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

	// init repo
	if opts.Init {
		// if init flag is set, no other options are expected
		if opts.Major || opts.Minor || opts.Patch || opts.Tag != "" || opts.Delete || opts.InitWithNoPrefix || opts.PreReleaseName != "" {
			fmt.Println("[*] Error: no arguments expected for command 'init'")
			fmt.Println("[*] Usage: taggo --init")
			os.Exit(1)
		} else {
			noPrefix := false
			err := repo.InitRepo(noPrefix)
			checkError(err)
			os.Exit(0)
		}

	}

	// init repo without prefix
	if opts.InitWithNoPrefix {
		// if init flag is set, no other options are expected
		if opts.Major || opts.Minor || opts.Patch || opts.Tag != "" || opts.Delete || opts.Init || opts.PreReleaseName != "" {
			fmt.Println("[*] Error: no arguments expected for command 'init-no-prefix'")
			fmt.Println("[*] Usage: taggo --init-no-prefix")
			os.Exit(1)
		} else {
			noPrefix := true
			err := repo.InitRepo(noPrefix)
			checkError(err)
			os.Exit(0)
		}
	}

	err := repo.Prerequisites()
	checkError(err)

	tag, err := repo.CurrentTag()
	checkError(err)

	fmt.Println("[*] Current tag: " + tag)

	if opts.Tag != "" && (opts.Major || opts.Minor || opts.Patch) {
		fmt.Println("[*] Error: tag and increment options are mutually exclusive")
		os.Exit(1)
	}

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
	if opts.PreReleaseName != "" {
		newTag, err = repo.CreatePreRelease(opts.PreReleaseName)
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
