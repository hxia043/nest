package image

import "github.com/spf13/cobra"

func NewAddImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Add image",
		Long:  "Add image",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	return cmd
}
