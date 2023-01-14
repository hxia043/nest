package image

import "github.com/spf13/cobra"

func NewDeleteImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete SUBCOMMAND",
		Short: "delete image or registry on nest",
		Long:  "delete image command will add the image to registry or add registry command will add registry info into the database of nest",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	return cmd
}
