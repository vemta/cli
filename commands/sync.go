package commands

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/gosuri/uilive"
	"github.com/spf13/cobra"
	"github.com/vemta/cli/cli"
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

func PullCommand(cmd *cobra.Command, args []string) {

	availableServices := make([]string, 0, len(cli.Services))
	for _, soft := range cli.Services {
		availableServices = append(availableServices, soft.Name)
	}

	res := []string{}

	prompt := &survey.MultiSelect{
		Message: "Select the repositories you want to pull/update",
		Options: availableServices,
	}

	survey.AskOne(prompt, &res)

	refresh := YesNoPrompt("Refresh services dependencies?", true)

	for _, source := range res {

		service, err := cli.GetServiceByName(source)
		if err != nil {
			panic(err)
		}

		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		repoDir := fmt.Sprintf("%s\\%s", dir, service.FolderName)

		writer := uilive.New()
		writer.Start()

		if _, err := os.Stat(dir + "/" + service.FolderName); os.IsNotExist(err) {
			fmt.Fprintf(writer, color.New(color.FgYellow).Sprintf("Cloning %s repository...\n", service.Name))
			if err := exec.Command("git", "clone", service.Repository).Run(); err != nil {
				fmt.Printf("Couldn't clone '%s': %s\n", service.Name, err.Error())
				continue
			}
			fmt.Fprintf(writer, color.New(color.FgGreen).Sprintf("Repository %s cloned successfully [✔]\n", service.Name))
		} else {
			fmt.Fprintf(writer, color.New(color.FgYellow).Sprintln("Synchronizing %s repository...\n"), service.Name)
			if err := exec.Command("git", "-C", repoDir, "reset", "--hard", "HEAD").Run(); err != nil {
				fmt.Fprintf(writer, color.New(color.FgRed).Sprintf("Couldn't clone %s repository [✘]\n", service.Name))
				continue
			}
			if err := exec.Command("git", "-C", repoDir, "pull", "origin", "master").Run(); err != nil {
				fmt.Fprintf(writer, color.New(color.FgRed).Sprintf("Couldn't clone %s repository [✘]\n", service.Name))
				continue
			}
			fmt.Fprintf(writer, color.New(color.FgGreen).Sprintf("Repository %s updated successfully [✔]\n", service.Name))
			writer.Stop()
		}

		if refresh {
			writer2 := uilive.New()
			writer2.Start()

			fmt.Fprintf(writer2, color.New(color.FgYellow).Sprintf("Downloading '%s' dependencies...\n", service.Name))

			cmd := exec.Command("go", "mod", "tidy")
			cmd.Dir = repoDir
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(writer2, color.New(color.FgRed).Sprintf("Couldn't download '%s' dependencies [✘]\n", service.Name))
				continue
			}

			fmt.Fprintf(writer2, color.New(color.FgGreen).Sprintf("Downloaded '%s' dependencies successfully [✔]\n", service.Name))
			writer2.Stop()
		}
		fmt.Println("")
	}

}
