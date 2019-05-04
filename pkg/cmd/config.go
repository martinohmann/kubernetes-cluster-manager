package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cmdutil"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/file"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type DumpConfigOptions struct {
	Filename string
	Output   string

	w io.Writer
}

func NewDumpConfigCommand(w io.Writer) *cobra.Command {
	o := &DumpConfigOptions{w: w}

	cmd := &cobra.Command{
		Use:   "dump-config [config-file]",
		Short: "Dumps the config to stdout",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("missing config-file argument")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd, args))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVar(&o.Output, "output", "yaml", "Output format")

	return cmd
}

func (o *DumpConfigOptions) Complete(cmd *cobra.Command, args []string) error {
	o.Filename = args[0]

	return nil
}

func (o *DumpConfigOptions) Validate() error {
	if o.Output != "" && o.Output != "yaml" && o.Output != "json" {
		return errors.New("--output must be 'yaml' or 'json'")
	}

	return nil
}

func (o *DumpConfigOptions) Run() error {
	opts := &Options{}

	if err := file.ReadYAML(o.Filename, opts); err != nil {
		return err
	}

	switch o.Output {
	case "json":
		buf, err := json.Marshal(opts)
		if err != nil {
			return err
		}

		fmt.Fprintln(o.w, string(buf))
	default:
		buf, err := yaml.Marshal(opts)
		if err != nil {
			return err
		}

		fmt.Fprint(o.w, string(buf))
	}

	return nil
}
