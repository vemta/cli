package commands

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

func LaunchCommand(cmd *cobra.Command, args []string) {
	res := []string{}

	prompt := &survey.MultiSelect{
		Message: "Select the repositories you want to pull/update",
		Options: []string{"mvc", "api", "payment"},
	}

	survey.AskOne(prompt, &res)
}
