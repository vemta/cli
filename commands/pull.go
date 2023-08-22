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
	res := []string{}

	prompt := &survey.MultiSelect{
		Message: "Select the repositories you want to pull/update",
		Options: []string{"mvc", "api", "payment"},
	}

	survey.AskOne(prompt, &res)

	refresh := YesNoPrompt("Refresh services dependencies?", true)

	for _, source := range res {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		repoDir := fmt.Sprintf("%s\\%s", dir, source)

		writer := uilive.New()
		writer.Start()

		if _, err := os.Stat(dir + "/" + source); os.IsNotExist(err) {
			fmt.Fprintf(writer, color.New(color.FgYellow).Sprintf("Cloning %s repository...\n", source))
			if err := exec.Command("git", "clone", "https://github.com/vemta/"+source).Run(); err != nil {
				fmt.Printf("Couldn't clone '%s': %s\n", source, err.Error())
				continue
			}
			fmt.Fprintf(writer, color.New(color.FgGreen).Sprintf("Repository %s cloned successfully [✔]\n", source))
		} else {
			fmt.Fprintf(writer, color.New(color.FgYellow).Sprintln("Synchronizing %s repository...\n"), source)
			if err := exec.Command("git", "-C", repoDir, "reset", "--hard", "HEAD").Run(); err != nil {
				fmt.Fprintf(writer, color.New(color.FgRed).Sprintf("Couldn't clone %s repository [✘]\n", source))
				continue
			}
			if err := exec.Command("git", "-C", repoDir, "pull", "origin", "master").Run(); err != nil {
				fmt.Fprintf(writer, color.New(color.FgRed).Sprintf("Couldn't clone %s repository [✘]\n", source))
				continue
			}
			fmt.Fprintf(writer, color.New(color.FgGreen).Sprintf("Repository %s updated successfully [✔]\n", source))
			writer.Stop()
		}

		if refresh {
			writer2 := uilive.New()
			writer2.Start()

			fmt.Fprintf(writer2, color.New(color.FgYellow).Sprintf("Downloading '%s' dependencies...\n", source))

			cmd := exec.Command("go", "mod", "tidy")
			cmd.Dir = repoDir
			if err := cmd.Run(); err != nil {
				fmt.Fprintf(writer2, color.New(color.FgRed).Sprintf("Couldn't download '%s' dependencies [✘]\n", source))
				continue
			}

			fmt.Fprintf(writer2, color.New(color.FgGreen).Sprintf("Downloaded '%s' dependencies successfully [✔]\n", source))
			writer2.Stop()
		}
		fmt.Println("")
	}

}
