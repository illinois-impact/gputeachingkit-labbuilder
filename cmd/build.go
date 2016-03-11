package cmd

import (
	"errors"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/abduld/wgx-pandoc/cmd/build"
)

var (
	buildOutputDir string
	showProgress   bool
	filterDocument bool
	formats        = []string{
		"pdf",
		"rtf",
		"html",
		"opendocument",
	}
	buildCmd = &cobra.Command{
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
			err := build.All("pdf", buildOutputDir, showProgress, filterDocument, args[0])
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
			err := build.All("markdown", buildOutputDir, showProgress, filterDocument, args[0])
			if err != nil {
				log.WithError(err).Error("Failed to generate Markdown labs")
				return err
			}
			return nil
		},
	}

	blackfridayBuildCmd = &cobra.Command{
		Use:     "blackfriday [./path/to/GPUTeachingKit-Labs] -o targetdir",
		Aliases: []string{},
		Short:   "Build the lab in Markdown format.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := build.All("blackfriday", buildOutputDir, showProgress, filterDocument, args[0])
			if err != nil {
				log.WithError(err).Error("Failed to generate Markdown labs")
				return err
			}
			return nil
		},
	}

	htmlBuildCmd = &cobra.Command{
		Use:     "html [./path/to/GPUTeachingKit-Labs] -o targetdir",
		Aliases: []string{"web", "html5"},
		Short:   "Build the lab in HTML format.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := build.All("html", buildOutputDir, showProgress, filterDocument, args[0])
			if err != nil {
				log.WithError(err).Error("Failed to generate HTML labs")
				return err
			}
			return nil
		},
	}

	rtfBuildCmd = &cobra.Command{
		Use:     "rtf [./path/to/GPUTeachingKit-Labs] -o targetdir",
		Aliases: []string{"doc"},
		Short:   "Build the lab in HTML format.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := build.All("rtf", buildOutputDir, showProgress, filterDocument, args[0])
			if err != nil {
				log.WithError(err).Error("Failed to generate HTML labs")
				return err
			}
			return nil
		},
	}

	opendocumentBuildCmd = &cobra.Command{
		Use:     "opendocument [./path/to/GPUTeachingKit-Labs] -o targetdir",
		Aliases: []string{"opendoc", "odf"},
		Short:   "Build the lab in HTML format.",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := build.All("opendocument", buildOutputDir, showProgress, filterDocument, args[0])
			if err != nil {
				log.WithError(err).Error("Failed to generate HTML labs")
				return err
			}
			return nil
		},
	}
	allBuildCmd = &cobra.Command{
		Use:     "all [./path/to/GPUTeachingKit-Labs] -o targetdir",
		Aliases: []string{"opendoc", "odf"},
		Short:   "Build the lab in HTML format.",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, format := range formats {
				err := build.All(format, buildOutputDir, showProgress, filterDocument, args[0])
				if err != nil {
					log.WithError(err).Error("Failed to generate HTML labs")
					return err
				}
			}
			return nil
		},
	}
)

func init() {
	buildCmd.PersistentFlags().StringVarP(&buildOutputDir, "output", "o", os.TempDir(), "The location of the output files.")
	buildCmd.PersistentFlags().BoolVarP(&showProgress, "progress", "p", false, "Prints the progress if enabled.")
	buildCmd.PersistentFlags().BoolVarP(&filterDocument, "filter", "f", true, "Pass the document through the pandoc filters.")
	buildCmd.AddCommand(allBuildCmd, pdfBuildCmd, markdownBuildCmd, blackfridayBuildCmd,
		htmlBuildCmd, rtfBuildCmd, opendocumentBuildCmd)
	RootCmd.AddCommand(buildCmd)
}
