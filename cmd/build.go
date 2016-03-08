package cmd

import (
	"errors"
	"os"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"gitlab.com/abduld/wgx-pandoc/cmd/build"
)

var (
	buildOutputDir string
	showProgress   bool
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
			err := build.All("pdf", buildOutputDir, showProgress, args[0])
			if err != nil {
				log.WithError(err).Error("Failed to generate PDF labs")
				return err
			}
			return nil
		},
	}

	markdownBuildCmd = &cobra.Command{
		Use:     "markdown [./path/to/GPUTeachingKit-Labs] -o targetdir",
		Aliases: []string{"md"},
		Short:   "Build the lab in Markdown format.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := build.All("markdown", buildOutputDir, showProgress, args[0])
			if err != nil {
				log.WithError(err).Error("Failed to generate Markdown labs")
				return err
			}
			return nil
		},
	}

	htmlBuildCmd = &cobra.Command{
		Use:     "html [./path/to/GPUTeachingKit-Labs] -o targetdir",
		Aliases: []string{"web"},
		Short:   "Build the lab in HTML format.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := build.All("html", buildOutputDir, showProgress, args[0])
			if err != nil {
				log.WithError(err).Error("Failed to generate HTML labs")
				return err
			}
			return nil
		},
	}
)

func init() {
	buildCmd.PersistentFlags().StringVarP(&buildOutputDir, "output", "o", os.TempDir(), "The location of the output files.")
	buildCmd.PersistentFlags().BoolVarP(&showProgress, "progress", "p", false, "Prints the progress if enabled.")
	buildCmd.AddCommand(pdfBuildCmd, markdownBuildCmd, htmlBuildCmd)
	RootCmd.AddCommand(buildCmd)
}
