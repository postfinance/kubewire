package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/postfinance/kubewire/pkg/report"
	"github.com/postfinance/kubewire/pkg/retrieval"
	"github.com/spf13/cobra"
)

// resourcesCmd represents the apiresources command
var resourcesCmd = &cobra.Command{
	Use:   "resources",
	Short: "List API resources",
	Long:  `List all API resources of a Kubernetes cluster including CustomResourceDefinitions.`,

	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := GetConfig()
		if err != nil {
			log.Fatalln(err)
		}

		data, err := retrieval.Resources(clientset)
		if err != nil {
			log.Fatal(err)
		}

		switch cmd.Flag("output").Value.String() {
		case "wide":
			printResourcesWide(data)
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
	rootCmd.AddCommand(resourcesCmd)

	resourcesCmd.Flags().StringP("output", "o", "wide", "Output format: json|yaml|wide")
}

func printResourcesWide(data []report.Resource) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "GroupVersion\tKind\tName\tNamespaced\tVerbs")

	for _, d := range data {
		fmt.Fprintf(w, "%s\t%s\t%s\t%t\t%s\n", d.GroupVersion, d.Kind, d.Name, d.Namespaced, d.Verbs)
	}

	w.Flush()
}
