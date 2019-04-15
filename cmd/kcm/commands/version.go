package commands

import (
	"encoding/json"
	"fmt"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/version"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var (
	shortVersion  bool
	versionOutput string

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Displays the version",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return validateOutput(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v := version.Get()

			if shortVersion {
				fmt.Printf("%s %s\n", rootCmd.Use, v.GitVersion)
				return nil
			}

			switch versionOutput {
			case "json":
				buf, err := json.Marshal(v)
				if err != nil {
					return err
				}

				fmt.Println(string(buf))
			case "yaml":
				buf, err := yaml.Marshal(v)
				if err != nil {
					return err
				}

				fmt.Println(string(buf))
			default:
				fmt.Printf("%s %#v\n", rootCmd.Use, v)
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&shortVersion, "short", false, "Display short version")
	versionCmd.Flags().StringVar(&versionOutput, "output", "", "Output format")
}
