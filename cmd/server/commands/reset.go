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
	resetCmd = &cobra.Command{
		Use:   "reset",
		Short: "Full database cleanup and deletion of user data",
		Run: func(cmd *cobra.Command, args []string) {

			app := BuildContainer(fx.Invoke(func(
				logger *logging.Logging,
				initialService *initial.Initial) {

				initialService.Reset()
			}))
			Start(app)
		},
	}
)
