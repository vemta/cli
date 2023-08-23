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

var repositories = &cobra.Command{
	Use:     "services",
	Aliases: []string{"svc"},
}

var pullCmd = &cobra.Command{
	Use:     "sync",
	Aliases: []string{"update", "pull", "synchronize"},
	Short:   "Pull and synchronize all the Vemta repositories",
	Run:     cli.SyncCommand,
}

func init() {

	repositories.AddCommand(pullCmd)
	root.AddCommand(repositories)
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
