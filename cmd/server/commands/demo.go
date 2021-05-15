package commands

import (
	. "github.com/e154/smart-home/cmd/server/container"
	. "github.com/e154/smart-home/common"
	"github.com/e154/smart-home/system/initial"
	"github.com/e154/smart-home/system/logging"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var (
	demoCmd = &cobra.Command{
		Use:   "demo",
		Short: "Filling the database with entities that help you quickly understand the process of setting up and managing the smart-home complex",
		Run: func(cmd *cobra.Command, args []string) {

			app := BuildContainer(fx.Invoke(func(
				logger *logging.Logging,
				initialService *initial.Initial) {

				initialService.InstallDemoData()
			}))
			Start(app)
		},
	}
)
