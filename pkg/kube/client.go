package kube

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/pytimer/k8sutil/apply"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	Clientset *kubernetes.Clientset
	config    *rest.Config
}

type ClientOption func(*Client) error

func WithConfig(config string) ClientOption {
	return func(c *Client) error {
		if config == "" {
			homeConfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
			log.Infof("try config file: %s", homeConfig)
			if _, err := os.Stat(homeConfig); err == nil {
				log.Infof("k8s using config: %s", homeConfig)
				config = homeConfig
			}
		}

		clientConfig, err := clientcmd.BuildConfigFromFlags("", config)
		if err != nil {
			return errors.Wrap(err, "failed to build config")
		}

		c.config = clientConfig
		return nil
	}
}

func Connect(opts ...ClientOption) (*Client, error) {
	client := &Client{}
	for _, opt := range opts {
		if err := opt(client); err != nil {
			return nil, err
		}
	}

	if client.config == nil {
		if err := WithConfig("")(client); err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(client.config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create client")
	}

	client.Clientset = clientset

	return client, errors.Wrapf(err, "failed to connect to k8s")
}

func (c *Client) Apply(template []byte) error {
	apply := apply.NewApplyOptions(dynamic.New(c.Clientset.RESTClient()), c.Clientset.Discovery())
	if err := apply.Apply(context.Background(), template); err != nil {
		return errors.Wrap(err, "apply config failed")
	}
	return nil
}

func (c *Client) CoreV1() v1.CoreV1Interface {
	return c.Clientset.CoreV1()
}
