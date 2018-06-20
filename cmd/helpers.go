package cmd

import (
	"encoding/json"
	"os"

	"github.com/postfinance/kubewire/pkg/access"
	yaml "gopkg.in/yaml.v2"
	"k8s.io/client-go/rest"
)

func GetConfig() (*rest.Config, error) {
	kcfg, err := rootCmd.PersistentFlags().GetString("kubeconfig")
	if err != nil {
		panic(err) // This should never occure
	}

	if kcfg == "" && access.IsInCluster() {
		return access.InCluster()
	}

	if kcfg != "" {
		return access.ForConfig(kcfg)
	}

	return access.Default()
}

func printJson(data interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(data)
}

func printYaml(data interface{}) {
	enc := yaml.NewEncoder(os.Stdout)
	enc.Encode(data)
}
