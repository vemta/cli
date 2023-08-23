package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
)

func YesNoPrompt(label string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}

func SyncCommand(cmd *cobra.Command, args []string) {

	availableServices := make([]string, 0, len(Services))
	for _, soft := range Services {
		availableServices = append(availableServices, soft.Name)
	}

	res := []string{}

	prompt := &survey.MultiSelect{
		Message: "Select the repositories you want to pull/update",
		Options: availableServices,
	}

	survey.AskOne(prompt, &res)

	for _, source := range res {

		service, err := GetServiceByName(source)
		if err != nil {
			panic(err)
		}

		workdir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		writer := uilive.New()
		writer.Start()

		if _, err := os.Stat(workdir + "/" + service.FolderName); os.IsNotExist(err) {
			fmt.Fprint(writer, processingMessage(fmt.Sprintf("Cloning %s repository...\n", service.Name)))
			if err := service.Clone(workdir); err != nil {
				fmt.Fprint(writer, errorMessage(fmt.Sprintf("Couldn't clone %s repository [✘]: %s\n", service.Name, err.Error())))
				continue
			}
			fmt.Fprint(writer, successMessage(fmt.Sprintf("Repository %s cloned successfully [✔]\n", service.Name)))
		} else {
			fmt.Fprint(writer, processingMessage(fmt.Sprintf("Synchronizing %s repository...\n", service.Name)))
			if err := service.Sync(workdir); err != nil {
				fmt.Fprint(writer, errorMessage(fmt.Sprintf("Couldn't clone %s repository [✘]: %s\n", service.Name, err.Error())))
				continue
			}
			fmt.Fprint(writer, successMessage(fmt.Sprintf("Repository %s updated successfully [✔]\n", service.Name)))
			writer.Stop()
		}
		fmt.Println("")
	}
}

func LaunchCommand(cmd *cobra.Command, args []string) {

}

func errorMessage(msg string) string {
	return color.New(color.FgRed).Sprintln(msg)
}

func processingMessage(msg string) string {
	return color.New(color.FgYellow).Sprintln(msg)
}

func successMessage(msg string) string {
	return color.New(color.FgGreen).Sprintln(msg)
}
