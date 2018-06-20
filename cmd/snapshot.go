package cmd

import (
	"log"
	"strings"
	"time"

	"github.com/postfinance/kubewire/pkg/report"
	"github.com/postfinance/kubewire/pkg/retrieval"
	"github.com/spf13/cobra"
)

// snapshotCmd represents the snapshot command
var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Take a snapshot of cluster resources and objects",
	Long: `Takes a snapshot of cluster resources and objects
that may be used by 'diff' to compare cluster states.
The output is either in json or yaml so that it can be processed by
other tools and also allows it to be stored in a database.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Split namespaces flag
		namespaces := strings.Split(cmd.Flag("namespaces").Value.String(), ",")
		if len(namespaces) == 1 && namespaces[0] == "" {
			namespaces = []string{}
		}

		// Create report
		rep, err := GetReport(namespaces)
		if err != nil {
			log.Fatalln(err)
		}

		// Printing
		switch cmd.Flag("output").Value.String() {
		case "json":
			printJson(rep)
		case "yaml":
			printYaml(rep)
		default:
			log.Fatalf("Unknown output format %s", cmd.Flag("output").Value.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)

	snapshotCmd.Flags().StringP("output", "o", "yaml", "Output format: json|yaml")
	snapshotCmd.Flags().StringP("namespaces", "n", "default,kube-public,kube-system", "Namespaces to scrape, commaseparated")
}

// GetReport creates a report
func GetReport(namespaces []string) (*report.Report, error) {
	rep := &report.Report{}
	rep.ScanStart = time.Now()
	rep.Configuration.KubewireVersion = Version
	rep.Configuration.Namespaces = namespaces

	// Setup
	clientset, err := GetConfig()
	if err != nil {
		return nil, err
	}

	// Server
	data, err := retrieval.ServerVersion(clientset)
	if err != nil {
		return nil, err
	}
	rep.Server = data

	// Resources
	data2, err := retrieval.Resources(clientset)
	if err != nil {
		return nil, err
	}
	rep.Resources = data2

	// Resource objects
	data3, err := retrieval.ResourceObjects(clientset, namespaces)
	if err != nil {
		return nil, err
	}
	rep.ResourceObjects = data3

	rep.ScanEnd = time.Now()

	return rep, nil
}
