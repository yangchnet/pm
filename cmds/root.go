package cmds

import "github.com/spf13/cobra"

var RootCmd = cobra.Command{Use: "pm"}

func init() {
	RootCmd.AddCommand(GenerateCmd())
	RootCmd.AddCommand(GetCmd())
	RootCmd.AddCommand(PushCmd())
	RootCmd.AddCommand(PullCmd())
	RootCmd.AddCommand(InitCmd())
}
