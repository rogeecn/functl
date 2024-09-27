package gen

import (
	config_map_to_path "github.com/rogeecn/functl/cmd/gen/confg_map_to_path"
	"github.com/rogeecn/functl/cmd/gen/model"
	"github.com/rogeecn/functl/cmd/gen/service_to_hosts"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "gen utils",
	}

	cmd.AddCommand(
		model.Command(),
		service_to_hosts.Command(),
		config_map_to_path.Command(),
	)

	return cmd
}
