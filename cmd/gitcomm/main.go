package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	cli "github.com/jawher/mow.cli"
	"github.com/studyzy/gitcomm"
	"github.com/studyzy/gitcomm/version"
)

func main() {
	app := cli.App("gitcomm", "Automate git commit messaging\n"+
		"\nSource https://github.com/studyzy/gitcomm")
	app.Version("V version", version.BuildDetails())

	app.Spec = "[-v] [-Ast] | [-u]"

	var (
		// declare the -r flag as a boolean flag
		addFiles = app.BoolOpt("A addAll", false, "Adds, modifies, and removes index entries "+
			"to match the working tree. Evals `git add -A`")
		verbose = app.BoolOpt("v verbose", false, "Switch log output")
		show    = app.BoolOpt("s show", false, "Show last commit or not. "+
			"Evals `git show -s` in the end of execution")
		tag  = app.BoolOpt("t tag", false, "Create an annonated tag for the next logical version")
		undo = app.BoolOpt("u undo", false, "Revert last commit")
	)

	// Specify the action to execute when the app is invoked correctly
	app.Action = func() {
		if !*verbose {
			log.SetFlags(0)
			log.SetOutput(ioutil.Discard)
		}
		if !gitcomm.CheckIsGitDir() {
			fmt.Println("Current directory is not inside git worktree")
			os.Exit(1)
		}
		if *undo {
			if gitcomm.PromptConfirm("Revert last commit?") {
				gitcomm.UndoLastCommit()
			}
			os.Exit(0)
		}
		if gitcomm.CheckForUncommited() {
			log.Printf("there are new changes in working directory\n")
			msg := gitcomm.Prompt()
			gitcomm.GitExec(*addFiles, *show, msg)
		} else {
			log.Printf("nothing to commit, working tree clean\n")
		}

		if *tag {
			level := gitcomm.TagPrompt()
			gitcomm.AutoTag(level)
		}
	}

	// Invoke the app passing in os.Args
	app.Run(os.Args)
}
