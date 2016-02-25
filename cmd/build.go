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
		Use:   "build [type] [./path/to/GPUTeachingKit-Labs] -o targetdir",
		Short: "Builds the lab dpeneding on the type",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("you must provide a path to the base teaching kit directory")
			}
			return nil
		},
	}

	pdfBuildCmd = &cobra.Command{
		Use:     "pdf [./path/to/GPUTeachingKit-Labs] -o targetdir",
		Aliases: []string{"PDF"},
		Short:   "Build the lab in PDF format.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return build.All("pdf", buildOutputDir, args[0])
		},
	}

	markdownBuildCmd = &cobra.Command{
		Use:     "markdown [./path/to/GPUTeachingKit-Labs] -o targetdir",
		Aliases: []string{"md"},
		Short:   "Build the lab in Markdown format.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return build.All("markdown", buildOutputDir, args[0])
		},
	}

	htmlBuildCmd = &cobra.Command{
		Use:     "html [./path/to/GPUTeachingKit-Labs] -o targetdir",
		Aliases: []string{"web"},
		Short:   "Build the lab in HTML format.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return build.All("html", buildOutputDir, args[0])
		},
	}
)

func init() {
	buildCmd.PersistentFlags().StringVarP(&buildOutputDir, "output", "o", os.TempDir(), "The location of the output files.")
	buildCmd.AddCommand(pdfBuildCmd, markdownBuildCmd, htmlBuildCmd)
	RootCmd.AddCommand(buildCmd)
}
