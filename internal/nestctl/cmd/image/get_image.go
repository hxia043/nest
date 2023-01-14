package image

import "github.com/spf13/cobra"

func NewGetImageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get SUBCOMMAND",
		Short: "get image or registry on nest",
		Long:  "get image command will add the image to registry or add registry command will add registry info into the database of nest",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	return cmd
}
