package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"functl/pkg/kube"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func genServiceToHostCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "service-to-host",
		Aliases: []string{"sth"},
		Short:   "gen k8s services to hosts file",
		RunE: func(cmd *cobra.Command, args []string) error {
			markStart := "#--------------------[start k8s svc hosts]--------------------"
			markEnd := "#--------------------[end k8s svc hosts]--------------------"

			hostFile := "/etc/hosts"

			// write to hosts file
			hosts, err := os.ReadFile(hostFile)
			if err != nil {
				return errors.Wrap(err, "failed to read hosts file")
			}

			lines := strings.Split(string(hosts), "\n")

			newLines := []string{}
			jumpLine := false
			for _, line := range lines {
				if line == markStart {
					jumpLine = true
				}

				if line == markEnd {
					jumpLine = false
				}

				if jumpLine {
					continue
				}

				newLines = append(newLines, line)
			}
			newLines = append(newLines, markStart)

			// connect to k8s
			client, err := kube.Connect()
			if err != nil {
				return err
			}
			// get all services
			services, err := client.CoreV1().Services("").List(context.Background(), metav1.ListOptions{})
			if err != nil {
				return errors.Wrap(err, "failed to list services")
			}

			servicesMap := map[string][]string{}
			for _, svc := range services.Items {
				if _, ok := servicesMap[svc.Spec.ClusterIP]; !ok {
					servicesMap[svc.Spec.ClusterIP] = []string{}
				}

				if svc.Spec.ClusterIP == "" || svc.Spec.ClusterIP == "None" {
					continue
				}

				name := svc.Name

				items := []string{
					svc.Spec.ClusterIP,
					name,
					fmt.Sprintf("%s.%s", name, svc.Namespace),
					fmt.Sprintf("%s.%s.svc", name, svc.Namespace),
					fmt.Sprintf("%s.%s.svc.cluster.local", name, svc.Namespace),
				}

				line := strings.Join(items, " ")
				fmt.Println(line)
				newLines = append(newLines, line)
			}
			newLines = append(newLines, markEnd)

			return os.WriteFile(hostFile, []byte(strings.Join(newLines, "\n")), 0o644)
		},
	}

	return cmd
}
