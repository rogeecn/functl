package gen

import "github.com/spf13/cobra"

func GenFile() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file",
		Short: "gen files",
	}

	return cmd
}

func genFile() {
}
