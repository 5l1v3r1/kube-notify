package main

import (
	"fmt"
	l "log"
	"os"
	"path/filepath"
	"time"

	"github.com/masahiro331/kube-notify/pkg/config"
	"github.com/masahiro331/kube-notify/pkg/signals"
	"github.com/masahiro331/kube-notify/pkg/slack"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeNotifyConf config.KubeNotifyConf

func main() {
	conf, err := config.Load("./config.toml")
	if err != nil {
		l.Fatal(err)
	}
	slack.Init(conf.Slack)

	config, err := getConfig(conf.KubeNotify)
	if err != nil {
		l.Fatal(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		l.Fatal(err.Error())
	}

	stopCh := signals.SetupSignalHandler()
	informerFactory := informers.NewSharedInformerFactory(clientset, time.Second*5)

	controller := NewController(clientset, informerFactory.Batch().V1().Jobs())
	informerFactory.Start(stopCh)

	if err = controller.Run(1, stopCh); err != nil {
		l.Fatalf("Error running controller: %s", err.Error())
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

func getConfig(conf config.KubeNotifyConf) (*rest.Config, error) {
	if conf.LocalMode {
		config, err := getOutClusterConfig(conf.ConfigPath)
		return config, err
	}
	config, err := getInClusterConfig()
	return config, err
}

func getInClusterConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}

func getOutClusterConfig(configPath string) (*rest.Config, error) {
	if configPath == "" {
		if home := homeDir(); home != "" {
			fmt.Printf("hogej:")
			configPath = filepath.Join(home, ".kube", "config")
		}
	}

	config, err := clientcmd.BuildConfigFromFlags("", configPath)
	if err != nil {
		return nil, err
	}

	return config, nil
}
