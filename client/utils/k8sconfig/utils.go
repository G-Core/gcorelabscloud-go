package k8sconfig

import (
	"os"

	"github.com/imdario/mergo"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/G-Core/gcorelabscloud-go/client/utils"
)

var defaultK8sConfigPath = "~/.kube/config"

func findK8sConfig(filename string) (string, error) {
	if filename == "" {
		filename = defaultK8sConfigPath
		env := os.Getenv("KUBECONFIG")
		if env != "" {
			filename = env
		}
	}
	return utils.GetAbsPath(filename)
}

func MergeKubeconfigFile(filename string, content []byte) (err error) {
	configPath, err := findK8sConfig(filename)
	if err != nil {
		return
	}
	config, err := clientcmd.LoadFromFile(configPath)
	if err != nil {
		return
	}
	newConfig, err := clientcmd.Load(content)
	if err != nil {
		return
	}
	err = mergo.Merge(config, newConfig, func(config2 *mergo.Config) {
		config2.Overwrite = true
	})
	if err != nil {
		return
	}
	return clientcmd.WriteToFile(*config, configPath)
}

func WriteKubeconfigFile(filename string, content []byte) error {
	configPath, err := findK8sConfig(filename)
	if err != nil {
		return err
	}
	config, err := clientcmd.Load(content)
	if err != nil {
		return err
	}
	return clientcmd.WriteToFile(*config, configPath)
}
