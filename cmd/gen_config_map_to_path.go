package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"functl/pkg/kube"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func genConfigMapToPath() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "config-map-to-path",
		Aliases: []string{"cmtp"},
		Short:   "gen k8s config map to path",
		Example: `functl gen config-map-to-path <namespace>/<configmap> <to-path>`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return errors.New("invalid args")
			}

			config := args[0]
			configs := strings.Split(config, "/")

			if len(configs) != 2 {
				return errors.New("invalid config <namespace>/<config-key>")
			}

			namespace, configMap := configs[0], configs[1]
			toPath := args[1]

			client, err := kube.Connect()
			if err != nil {
				return err
			}
			// get configmap
			cm, err := client.CoreV1().ConfigMaps(namespace).Get(context.Background(), configMap, metav1.GetOptions{})
			if err != nil {
				return errors.Wrap(err, "failed to list services")
			}

			for name, content := range cm.Data {
				fmt.Println(name, "------------------------")
				fmt.Println(content)
				if err := os.WriteFile(filepath.Join(toPath, name), []byte(content), os.ModePerm); err != nil {
					return errors.Wrap(err, "failed to write file")
				}
			}
			return nil
		},
	}

	return cmd
}
