package cmd

import "github.com/spf13/cobra"

func GenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gen",
		Short: "gen utils",
	}

	cmd.AddCommand(
		genModelCommand(),
	)

	return cmd
}
