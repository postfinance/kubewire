package access

import (
	"os"

	// These are needed in order to support the authentication methods
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	"k8s.io/client-go/rest"

	"k8s.io/client-go/tools/clientcmd"
)

// IsInCluster returns true if the process is running in a Pod
func IsInCluster() bool {
	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	port := os.Getenv("KUBERNETES_SERVICE_PORT")
	return len(host) != 0 && len(port) != 0

}

// InCluster returns a REST configuration for using it inside a Pod
func InCluster() (*rest.Config, error) {
	return rest.InClusterConfig()
}

// ForConfig returns a client config from the provided file path
func ForConfig(kcfg string) (*rest.Config, error) {
	fileConfig, err := clientcmd.LoadFromFile(kcfg)
	if err != nil {
		return nil, err
	}

	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewDefaultClientConfig(*fileConfig, configOverrides)

	return kubeConfig.ClientConfig()
}

// Default returns the client config using the default kubeconfig search paths
func Default() (*rest.Config, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	return kubeConfig.ClientConfig()
}
