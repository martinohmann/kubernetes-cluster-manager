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
		Use:   "dump-config",
		Short: "Dumps the config to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Complete(cmd))
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.Flags().StringVar(&o.Output, "output", "", "Output format")
	cmdutil.AddConfigFlag(cmd)

	return cmd
}

func (o *DumpConfigOptions) Complete(cmd *cobra.Command) error {
	o.Filename = cmdutil.GetString(cmd, "config")

	return nil
}

func (o *DumpConfigOptions) Validate() error {
	if o.Output != "" && o.Output != "yaml" && o.Output != "json" {
		return errors.New("--output must be 'yaml' or 'json'")
	}

	if o.Filename != "" && !file.Exists(o.Filename) {
		return errors.Errorf("File %q does not exist", o.Filename)
	}

	return nil
}

func (o *DumpConfigOptions) Run() error {
	opts := &Options{}

	if err := opts.MergeConfig(o.Filename); err != nil {
		return err
	}

	switch o.Output {
	case "json":
		buf, err := json.Marshal(opts)
		if err != nil {
			return err
		}

		fmt.Fprintln(o.w, string(buf))
	case "yaml":
		buf, err := yaml.Marshal(opts)
		if err != nil {
			return err
		}

		fmt.Fprintln(o.w, string(buf))
	default:
		fmt.Fprintf(o.w, "%#v\n", opts)
	}

	return nil
}
