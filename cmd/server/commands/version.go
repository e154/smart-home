package commands

import (
	"fmt"
	"github.com/e154/smart-home/version"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "The current version of the smart-home software package",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(version.ShortVersionBanner, version.GetHumanVersion())
		},
	}
)
