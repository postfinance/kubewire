package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/tabwriter"

	"github.com/postfinance/kubewire/pkg/report"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare snapshots with another or a live cluster",
	Long: `Compares the state from the baseline with
the current state of the cluster or another snapshot. The namespaces defined in the
baseline report are used for snapshotting the live cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read report
		baselineFile := cmd.Flag("baseline").Value.String()
		raw, err := ioutil.ReadFile(baselineFile)
		if err != nil {
			log.Fatalln(err)
		}

		baseline := report.Report{}
		err = yaml.Unmarshal(raw, &baseline)
		if err != nil {
			log.Fatalln(err)
		}

		snapshotFile := cmd.Flag("snapshot").Value.String()
		live := &report.Report{}

		if snapshotFile == "" {
			// Create report
			var err error
			live, err = GetReport(baseline.Configuration.Namespaces)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			// Read snapshot
			raw, err := ioutil.ReadFile(snapshotFile)
			if err != nil {
				log.Fatalln(err)
			}

			err = yaml.Unmarshal(raw, live)
			if err != nil {
				log.Fatalln(err)
			}
		}

		// Diff
		data := report.DiffReports(baseline, *live)

		// Printing
		switch cmd.Flag("output").Value.String() {
		case "wide":
			printDiffWide(data)
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
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().StringP("baseline", "b", "baseline.yaml", "Baseline report in yaml format")
	diffCmd.Flags().StringP("snapshot", "s", "", "Snapshot in yaml format to read in, empty to run against live cluster")
	diffCmd.Flags().StringP("output", "o", "wide", "Output format: json|yaml|wide")
}

func printDiffWide(data []report.DiffReport) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "Element\tA\tB")

	for _, d := range data {
		fmt.Fprintf(w, "%s\t%s\t%s\n", d.Element, d.A, d.B)
	}

	w.Flush()
}
