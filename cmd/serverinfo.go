package cmd

import (
	"fmt"
	"log"

	"github.com/postfinance/kubewire/pkg/retrieval"
	"github.com/spf13/cobra"
)

// serverinfoCmd represents the serverinfo command
var serverinfoCmd = &cobra.Command{
	Use:   "serverinfo",
	Short: "Prints server info",
	Long:  `Prints informations about the remote server.`,
	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := GetConfig()
		if err != nil {
			log.Fatalln(err)
		}

		version, err := retrieval.ServerVersion(clientset)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("Host: %s, Version: %s\n", version.Host, version.Version)
	},
}

func init() {
	rootCmd.AddCommand(serverinfoCmd)
}
