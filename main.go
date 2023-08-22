package main

import (
	"log"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
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

	if _, err := exec.LookPath("git"); err != nil {
		log.Fatal(color.New(color.FgRed).Sprint("Couldn't find Git! Make sure it is installed and added to the path."))
		return
	}

	if _, err := exec.LookPath("docker"); err != nil {
		log.Fatal(color.New(color.FgRed).Sprint("Couldn't find Docker! Make sure it is installed and added to the path."))
		return
	}

	if _, err := exec.LookPath("docker-compose"); err != nil {
		log.Fatal(color.New(color.FgRed).Sprint("Couldn't find docker-compose! Make sure it is installed and added to the path."))
		return
	}

	root.Execute()

}
