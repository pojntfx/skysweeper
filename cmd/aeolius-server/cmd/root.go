package cmd

import (
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	postgresURLFlag = "postgres-url"
	laddrFlag       = "laddr"
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

		if v := os.Getenv("DATABASE_URL"); v != "" {
			log.Println("Using database address from DATABASE_URL env variable")

			viper.Set(postgresURLFlag, v)
		}

		if v := os.Getenv("PORT"); v != "" {
			log.Println("Using port from PORT env variable")

			la, err := net.ResolveTCPAddr("tcp", viper.GetString(laddrFlag))
			if err != nil {
				return err
			}

			p, err := strconv.Atoi(v)
			if err != nil {
				return err
			}

			la.Port = p
			viper.Set(laddrFlag, la.String())
		}

		return nil
	},
}

func Execute() error {
	rootCmd.PersistentFlags().String(postgresURLFlag, "postgresql://postgres@localhost:5432/aeolius?sslmode=disable", "PostgreSQL URL (can also be set using `DATABASE_URL` env variable)")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		return err
	}

	viper.AutomaticEnv()

	return rootCmd.Execute()
}
