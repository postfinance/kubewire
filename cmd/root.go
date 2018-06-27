package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version defines the kubewire version string
var Version = "devel"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kubewire",
	Short: "Integrity checker for Kubernetes",
	Long: `kubewire is a Kubernetes integrity checker which acts as a tripwire for global
Kubernetes or namespaced resources that could impact the
whole cluster.

It detects if it is running in a Kubernetes cluster and uses the service account
of the Pod if available. If this is not the case, it looks through the default kubectl
paths for a kubeconfig. Either case can be overriden by setting the 'kubeconfig' flag.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = Version

	rootCmd.PersistentFlags().StringP("kubeconfig", "k", "", "absolute path to the kubeconfig file")
}
