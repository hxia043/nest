package cmd

import "github.com/spf13/cobra"

func NewSWICCtlExecute() *cobra.Command {
	cmds := &cobra.Command{
		Use:   "swicctl",
		Short: "swic is a software image controller",
		Long:  "swic is a software image controller",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmds.CompletionOptions.DisableDefaultCmd = true

	cg := CommandGroups{
		{
			Message: "Basic Command:",
			Commands: []*cobra.Command{
				NewAddCmd(),
				NewDeleteCmd(),
				NewUpdateCmd(),
				NewGetCmd(),
				NewShowCmd(),
				NewOnboardCmd(),
			},
		},
	}

	cg.Add(cmds)

	return cmds
}
