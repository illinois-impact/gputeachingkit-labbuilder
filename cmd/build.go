package cmd

import (
	"errors"
	"os"

	"strings"

	log "github.com/rai-project/logger"
	"github.com/spf13/cobra"
	"bitbucket.org/hwuligans/gputeachingkit-labbuilder/cmd/build"
)

var (
	buildOutputDir string
	showProgress   bool
	filterDocument bool
	logger         = log.StandardLogger()
	formats        = map[string][]string{
		"pdf":         []string{"adobepdf"},
		"rtf":         []string{"richtext"},
		"html":        []string{"web"},
		"docx":        []string{"word"},
		"markdown":    []string{"md"},
		"blackfriday": []string{"basichtml"},
		"blackfridaytex": []string{
			"tex",
			"latex",
			"blackfriday_tex",
		},
		"opendocument": []string{"odf"},
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
)

func init() {

	buildCmd.PersistentFlags().StringVarP(&buildOutputDir, "output", "o", os.TempDir(), "The location of the output files.")
	buildCmd.PersistentFlags().BoolVarP(&showProgress, "progress", "p", false, "Prints the progress if enabled.")
	buildCmd.PersistentFlags().BoolVarP(&filterDocument, "filter", "f", true, "Pass the document through the pandoc filters.")

	for format0, aliases0 := range formats {
		format := format0
		aliases := aliases0
		buildCmd.AddCommand(&cobra.Command{
			Use:     format + " [./path/to/GPUTeachingKit-Labs] -o targetdir",
			Aliases: aliases,
			Short:   "Build the lab in " + strings.ToUpper(format) + " format.",
			RunE: func(cmd *cobra.Command, args []string) error {
				err := build.All(format, buildOutputDir, showProgress, filterDocument, args[0])
				if err != nil {
					log.WithError(err).Error("✖ Failed to generate " + strings.ToUpper(format) + " labs")
					return err
				}
				return nil
			},
		})
	}
	buildCmd.AddCommand(&cobra.Command{
		Use:   "all [./path/to/GPUTeachingKit-Labs] -o targetdir",
		Short: "Build the lab using all format.",
		RunE: func(cmd *cobra.Command, args []string) error {
			for format := range formats {
				err := build.All(format, buildOutputDir, showProgress, filterDocument, args[0])
				if err != nil {
					log.WithError(err).Error("✖ Failed to generate " + strings.ToUpper(format) + " labs")
					return err
				}
			}
			return nil
		},
	})
	RootCmd.AddCommand(buildCmd)
}
