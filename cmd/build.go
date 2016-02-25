package cmd

import (
	"errors"

	"os"

	"github.com/spf13/cobra"
	"gitlab.com/abduld/wgx-pandoc/cmd/build"
)

var (
	buildOutputDir string
	buildCmd       = &cobra.Command{
		Use:   "build [./path/to/GPUTeachingKit-Labs]",
		Short: "Makes the lab using the same mechanism as the make_lab_handout.py",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("you must provide a path to the base directory")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return build.All(buildOutputDir, args[0])
		},
	}
)

func init() {
	buildCmd.Flags().StringVarP(&buildOutputDir, "output", "o", os.TempDir(), "The location of the output files.")
	RootCmd.AddCommand(buildCmd)
}
