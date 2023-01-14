package cmd

import (
	"fmt"
	"os"

	"github.com/hxia043/nest/internal/swicctl/cmd/image"
	"github.com/hxia043/nest/internal/swicctl/cmd/registry"
	"github.com/spf13/cobra"
)

type CommandGroup struct {
	Message  string
	Commands []*cobra.Command
}

type CommandGroups []CommandGroup

func NewAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add image or registry resource",
		Long:  "Add image or registry resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(registry.NewAddRegistryCmd())
	cmd.AddCommand(image.NewAddImageCmd())

	return cmd
}

func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [registry|image]",
		Short: "Delete image or registry resource",
		Long:  "Delete image or registry resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(registry.NewDeleteRegistryCmd())
	cmd.AddCommand(image.NewDeleteImageCmd())

	return cmd
}

func NewUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [registry]",
		Short: "Update registry resource",
		Long:  "Update registry resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(registry.NewUpdateRegistryCmd())
	cmd.AddCommand(image.NewUpdateImageCmd())

	return cmd
}

func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [image|registry]",
		Short: "Get image or registry resource",
		Long:  "Get image or registry resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(registry.NewGetRegistryCmd())
	cmd.AddCommand(image.NewGetImageCmd())

	return cmd
}

func NewShowCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [image|registry]",
		Short: "Show image or registry resource",
		Long:  "Show image or registry resource",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(registry.NewShowRegistryCmd())
	cmd.AddCommand(image.NewGetImageCmd())

	return cmd
}

func NewOnboardCmd() *cobra.Command {
	o := image.ImageOnboardOption{}
	cmd := &cobra.Command{
		Use:     "onboard [name]",
		Short:   "Onboard images",
		Long:    "Onboard images",
		Example: "swicctl onboard test -i /var/home/core/images/ -n cran1",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				return
			}

			if err := o.Onboard(args); err != nil {
				fmt.Fprint(os.Stdout, err.Error()+"\n")
			}
		},
	}

	cmd.Flags().StringVarP(&o.ImagePath, "image-path", "i", "", "local directory of images (default \"current directory\")")
	cmd.Flags().StringVarP(&o.Namespace, "namespace", "n", "", "onboard namespace of registry")

	return cmd
}

func (g CommandGroups) Add(c *cobra.Command) {
	for _, group := range g {
		c.AddCommand(group.Commands...)
	}
}
