package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/martinohmann/kubernetes-cluster-manager/pkg/cmdutil"
	"github.com/martinohmann/kubernetes-cluster-manager/pkg/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

type VersionOptions struct {
	Short  bool
	Output string

	w io.Writer
}

func NewVersionCommand(w io.Writer) *cobra.Command {
	o := &VersionOptions{w: w}

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Displays the version",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Validate())
			cmdutil.CheckErr(o.Run())
		},
	}

	cmd.Flags().BoolVar(&o.Short, "short", false, "Display short version")
	cmd.Flags().StringVar(&o.Output, "output", "", "Output format")

	return cmd
}

func (o *VersionOptions) Validate() error {
	if o.Output != "" && o.Output != "yaml" && o.Output != "json" {
		return errors.New("--output must be 'yaml' or 'json'")
	}

	return nil
}

func (o *VersionOptions) Run() error {
	v := version.Get()

	if o.Short {
		fmt.Fprintln(o.w, v.GitVersion)
		return nil
	}

	switch o.Output {
	case "json":
		buf, err := json.Marshal(v)
		if err != nil {
			return err
		}

		fmt.Fprintln(o.w, string(buf))
	case "yaml":
		buf, err := yaml.Marshal(v)
		if err != nil {
			return err
		}

		fmt.Fprintln(o.w, string(buf))
	default:
		fmt.Fprintf(o.w, "%#v\n", v)
	}

	return nil
}
