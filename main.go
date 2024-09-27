/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	_ "embed"
	"log"
	"os"
	"path/filepath"

	"functl/cmd/gen"
	"functl/config"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

//go:embed fun_tpl.yaml
var defaultConfig []byte
var configFile string

func init() {
	configFile = filepath.Join(lo.Must(os.Getwd()), ".fun.yaml")
	// if configFile not exists, create it
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		os.WriteFile(configFile, defaultConfig, 0o644)
		log.Fatal("config file not found, created a default one, please fill it")
	}
}

func main() {
	rootCmd := &cobra.Command{
		SilenceErrors: true,
		SilenceUsage:  true,
		Use:           "functl",
		Short:         "A brief description of your application",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return config.Load(configFile)
		},
	}

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Println(err)
		cmd.Println(cmd.UsageString())
		return err
	})

	rootCmd.AddCommand(
		gen.Command(),
	)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
