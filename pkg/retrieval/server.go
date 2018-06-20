package retrieval

import (
	"github.com/postfinance/kubewire/pkg/report"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ServerVersion retrieves the server version of the remote server
func ServerVersion(config *rest.Config) (report.Server, error) {
	rep := report.Server{
		Host: config.Host,
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return rep, err
	}

	version, err := clientset.ServerVersion()
	if err != nil {
		return rep, err
	}

	rep.Version = version.String()

	return rep, nil
}
