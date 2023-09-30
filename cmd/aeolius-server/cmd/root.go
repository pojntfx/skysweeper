package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	postgresURLFlag = "postgres-url"
)

var rootCmd = &cobra.Command{
	Use:   "aeolius-server",
	Short: "Start Aeolius managers and workers",
	Long: `Automatically delete your old skeets from Bluesky.
Find more information at:
https://github.com/pojntfx/aeolius`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		viper.SetEnvPrefix("aeolius")
		viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}

		return nil
	},
}

func Execute() error {
	rootCmd.PersistentFlags().String(postgresURLFlag, "postgresql://postgres@localhost:5432/aeolius?sslmode=disable", "PostgreSQL URL")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		return err
	}

	viper.AutomaticEnv()

	return rootCmd.Execute()
}
