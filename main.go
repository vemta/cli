package main

import (
	"github.com/spf13/cobra"
	"github.com/vemta/cli/commands"
)

var root = &cobra.Command{
	Use:   "vemta",
	Short: "Vemta microservices management CLI",
}

var pullCmd = &cobra.Command{
	Use:     "pull",
	Aliases: []string{"update"},
	Short:   "Pull and update all the Vemta repositories",
	Run:     commands.PullCommand,
}

func init() {
	root.AddCommand(pullCmd)
}

func main() {

	root.Execute()

}
