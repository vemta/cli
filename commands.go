package cli

import (
	"bufio"
	"context"
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
	ctx := context.Background()
	availableServices := make([]string, 0, len(Services))
	for _, soft := range Services {
		availableServices = append(availableServices, soft.Name)
	}

	res := []string{}
	restart, err := cmd.Flags().GetBool("restart")
	if err != nil {
		panic(err)
	}

	prompt := &survey.MultiSelect{
		Message: "Select the services you want to launch/restart",
		Options: availableServices,
	}

	survey.AskOne(prompt, &res)

	for _, source := range res {

		service, err := GetServiceByName(source)
		if err != nil {
			panic(err)
		}

		parentWriter := uilive.New()
		parentWriter.Start()

		containers := GetContainersOfService(service)

		if len(*containers) <= 0 {
			fmt.Fprint(parentWriter, errorMessage(fmt.Sprintf("No containers found for service %s. Make sure you have executed the command: vemta services setup.\n", service.Name)))
			continue
		}

		fmt.Fprint(parentWriter, processingMessage(fmt.Sprintf("↑ Launching service %s...[0/%d]", service.Name, len(*containers))))

		failedCount := 0
		successCount := 0

		for _, container := range *containers {

			containerWriter := uilive.New()
			containerWriter.Start()

			fmt.Fprint(containerWriter, processingMessage(fmt.Sprintf("    - Launching container %s", container.Name)))

			running, er := IsContainerRunning(ctx, &container)
			if er != nil {
				failedCount++
				fmt.Fprint(containerWriter, errorMessage(fmt.Sprintf("    ✘ Container %s launch failed\n", container.Name)))
				fmt.Printf(errorMessage(fmt.Sprintf("        ✘ Coudln't retrieve container stats. Make sure you have execute the command: vemta services setup\n")))
				continue
			}
			fmt.Printf(successMessage(fmt.Sprintf("        ✔ Container stats retrieved successfully!\n")))
			if running {
				if restart {
					stoppingWriter := uilive.New()
					stoppingWriter.Start()
					fmt.Fprint(stoppingWriter, processingMessage("        - Stopping container..."))
					if err := StopContainer(ctx, &container); err != nil {
						fmt.Fprint(containerWriter, errorMessage(fmt.Sprintf("    ✘ Container %s launch failed\n", container.Name)))
						fmt.Fprint(stoppingWriter, errorMessage(fmt.Sprintf("        ✘ Couldn't stop container: %s\n", err.Error())))
						stoppingWriter.Stop()
						failedCount++
						continue
					}
					fmt.Fprint(stoppingWriter, successMessage(fmt.Sprintf("        ✔ Container stopped successfully!\n")))
					stoppingWriter.Stop()
				} else {
					fmt.Printf(successMessage(fmt.Sprintf("        ✔ Container already launched!\n")))
					fmt.Fprint(containerWriter, processingMessage(fmt.Sprintf("    Container %s launched successfully\n", container.Name)))
					successCount++
					fmt.Fprint(parentWriter, processingMessage(fmt.Sprintf("↑ Launching service %s... [%d/%d]", service.Name, successCount, (*containers))))
					continue
				}
			}
			if !running || (running && restart) {
				launchingWriter := uilive.New()
				launchingWriter.Start()
				fmt.Fprint(launchingWriter, processingMessage("        - Starting container..."))
				if err := LaunchContainer(ctx, &container); err != nil {
					fmt.Fprint(containerWriter, errorMessage(fmt.Sprintf("    ✘ Container %s launch failed\n", container.Name)))
					fmt.Fprint(launchingWriter, errorMessage(fmt.Sprintf("        ✘ Couldn't launch container: %s\n", err.Error())))
					launchingWriter.Stop()
					failedCount++
					continue
				}
				fmt.Fprint(launchingWriter, successMessage(fmt.Sprintf("        ✔ Container launched successfully!\n")))
				successCount++
				fmt.Fprint(parentWriter, processingMessage(fmt.Sprintf("↑ Launching service %s... [%d/%d]", service.Name, successCount, (*containers))))
				launchingWriter.Stop()
			}
			containerWriter.Stop()

		}

		parentWriter.Stop()
		fmt.Println("")
	}
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
