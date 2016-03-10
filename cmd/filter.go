package cmd

import (
	"errors"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/abduld/wgx-pandoc/cmd/filter"
)

var (
	filterOutputFile string
	filterFormat     string
	filterCmd        = &cobra.Command{
		Use:   "filter [./path/to/markdownFile] -o targetFile",
		Short: "Runs the markdown through filters to mutate the document",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("you must provide a path to the input markdown file")
			}
			_, err := filter.Filter(filterOutputFile, args[0], filterFormat)
			if err != nil {
				log.WithError(err).Error("Failed to filter markdown file.")
				return err
			}
			return nil
		},
	}
)

func init() {
	dir, err := os.Getwd()
	if err != nil {
		dir = os.TempDir()
	}
	filterCmd.PersistentFlags().StringVarP(&filterOutputFile, "output", "o", dir, "The directory of the output file.")
	filterCmd.PersistentFlags().StringVarP(&filterFormat, "format", "f", "", "The format filter.")
	RootCmd.AddCommand(filterCmd)
}
