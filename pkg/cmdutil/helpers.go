package cmdutil

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func GetBool(cmd *cobra.Command, flag string) bool {
	v, err := cmd.Flags().GetBool(flag)
	if err != nil {
		log.Fatal(err)
	}

	return v
}

func GetString(cmd *cobra.Command, flag string) string {
	v, err := cmd.Flags().GetString(flag)
	if err != nil {
		log.Fatal(err)
	}

	return v
}

func GetInt(cmd *cobra.Command, flag string) int {
	v, err := cmd.Flags().GetInt(flag)
	if err != nil {
		log.Fatal(err)
	}

	return v
}
