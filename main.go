package main

import (
	"log"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/vemta/cli/cli"
	"github.com/vemta/cli/commands"
)

var root = &cobra.Command{
	Use:   "vemta",
	Short: "Vemta microservices management CLI",
}

var repositories = &cobra.Command{
	Use:     "services",
	Aliases: []string{"svc"},
}

var pullCmd = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"update", "pull"},
	Short:   "Pull and synchronize all the Vemta repositories",
	Run:     commands.PullCommand,
}

func init() {
	repositories.AddCommand()
	root.AddCommand(repositories)
}

func main() {

	for _, software := range cli.MustHaveSoftwares {
		if _, err := exec.LookPath(software); err != nil {
			log.Fatal(color.New(color.FgRed).Sprintf("Couldn't find %s! Make sure it is installed and added to the path.", software))
			return
		}
	}

	root.Execute()

}
