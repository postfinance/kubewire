package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/postfinance/kubewire/pkg/report"
	"github.com/postfinance/kubewire/pkg/retrieval"
	"github.com/spf13/cobra"
)

// resourceobjectsCmd represents the apiresources command
var resourceobjectsCmd = &cobra.Command{
	Use:   "resourceobjects",
	Short: "List API resource objects",
	Long: `List all global resource objects of a Kubernetes cluster,
This also includes CustomResourceDefinitions. Resource objects
in namespaces that are considered to have a global impact are also listed.
These namespaces can be customized with the 'namespaces' flag.`,

	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := GetConfig()
		if err != nil {
			log.Fatalln(err)
		}

		namespaces := strings.Split(cmd.Flag("namespaces").Value.String(), ",")
		if len(namespaces) == 1 && namespaces[0] == "" {
			namespaces = []string{}
		}

		data, err := retrieval.ResourceObjects(clientset, namespaces)
		if err != nil {
			log.Fatal(err)
		}

		switch cmd.Flag("output").Value.String() {
		case "wide":
			printResourceObjectsWide(data)
		case "json":
			printJson(data)
		case "yaml":
			printYaml(data)
		default:
			log.Fatalf("Unknown output format %s", cmd.Flag("output").Value.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(resourceobjectsCmd)

	resourceobjectsCmd.Flags().StringP("output", "o", "wide", "Output format: json|yaml|wide")
	resourceobjectsCmd.Flags().StringP("namespaces", "n", "default,kube-public,kube-system", "Namespaces to scrape, commaseparated")
}

func printResourceObjectsWide(data []report.ResourceObject) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "GroupVersion\tResource\tNamespace\tName")

	for _, d := range data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", d.GroupVersion, d.Resource, d.Namespace, d.Name)
	}

	w.Flush()
}
