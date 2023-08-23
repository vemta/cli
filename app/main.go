package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vemta/cli"
)

var root = &cobra.Command{
	Use:   "vemta",
	Short: "Vemta microservices management CLI",
}

var servicesCmd = &cobra.Command{
	Use:     "services",
	Aliases: []string{"svc"},
}

var pullCmd = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"update", "pull", "synchronize"},
	Short:   "Pull and synchronize Vemta repositories",
	Run:     cli.SyncCommand,
}

var launchCmd = &cobra.Command{
	Use:     "launch",
	Aliases: []string{"up", "start"},
	Short:   "Launch Vemta services's containers",
	Run:     cli.LaunchCommand,
}

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"up"},
	Short:   "Build venta services",
	Run:     cli.BuildCommand,
}

func init() {
	launchCmd.Flags().BoolP("restart", "r", false, "Restart the container if it is already launched")
	servicesCmd.AddCommand(launchCmd)
	servicesCmd.AddCommand(pullCmd)
	servicesCmd.AddCommand(buildCmd)
	root.AddCommand(servicesCmd)
}

func main() {

	if err := viper.SafeWriteConfigAs("./config.json"); err != nil {
		if os.IsNotExist(err) {
			err = viper.WriteConfigAs("./config.json")
			if err != nil {
				panic(err)
			}
		}
	}

	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.ReadInConfig()

	/*for _, software := range cli.MustHaveSoftwares {
		if _, err := exec.LookPath(software); err != nil {
			log.Fatal(color.New(color.FgRed).Sprintf("Couldn't find %s! Make sure it is installed and added to the path.", software))
			return
		}
	}*/

	root.Execute()

}
