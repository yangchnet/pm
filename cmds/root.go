package cmds

import "github.com/spf13/cobra"

var RootCmd = cobra.Command{
	Use:     "pm",
	Version: "__VERSION__",
}

func init() {
	RootCmd.AddCommand(GenerateCmd())
	RootCmd.AddCommand(GetCmd())
	RootCmd.AddCommand(PushCmd())
	RootCmd.AddCommand(PullCmd())
	RootCmd.AddCommand(InitCmd())
	RootCmd.AddCommand(DelCmd())
	RootCmd.AddCommand(ListCmd())
	RootCmd.AddCommand(UpdateCmd())
}
