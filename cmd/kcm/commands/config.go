package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var (
	configDumpOutput string

	dumpConfigCmd = &cobra.Command{
		Use:   "dump-config",
		Short: "Dumps the config to stdout",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateOutput(cmd)
		},
		RunE: dumpConfig,
	}
)

func init() {
	rootCmd.AddCommand(dumpConfigCmd)
	dumpConfigCmd.Flags().StringVar(&configDumpOutput, "output", "", "Output format")
}

func dumpConfig(cmd *cobra.Command, args []string) error {
	switch configDumpOutput {
	case "json":
		buf, err := json.Marshal(cfg)
		if err != nil {
			return err
		}

		fmt.Println(string(buf))
	case "yaml":
		buf, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}

		fmt.Println(string(buf))
	default:
		fmt.Printf("%s config %#v\n", rootCmd.Use, cfg)
	}

	return nil
}
