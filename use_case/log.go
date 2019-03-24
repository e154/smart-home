package use_case

type LogCommand struct {
	*CommonCommand
}

func NewLogCommand(common *CommonCommand) *LogCommand {
	return &LogCommand{
		CommonCommand: common,
	}
}
