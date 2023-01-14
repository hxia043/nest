package registry

import (
	"fmt"
	"os"

	"github.com/hxia043/nest/internal/swicctl/db"
	"github.com/spf13/cobra"
)

type registryDeleteOption struct{}

func newRegistryDeleteOption() *registryDeleteOption {
	return &registryDeleteOption{}
}

func (o *registryDeleteOption) run(args []string) error {
	for _, name := range args {
		dbOption, err := db.InitDB()
		if err != nil {
			fmt.Fprint(os.Stdout, err.Error()+"\n")
		}

		if err := dbOption.Delete(name); err != nil {
			fmt.Fprint(os.Stdout, err.Error()+"\n")
		}
	}

	return nil
}

func NewDeleteRegistryCmd() *cobra.Command {
	o := newRegistryDeleteOption()
	cmd := &cobra.Command{
		Use:     "registry [name]",
		Short:   "Delete registry",
		Long:    "Delete registry",
		Example: "swic delete registry test",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				_ = cmd.Help()
				return
			}

			if err := o.run(args); err != nil {
				fmt.Fprint(os.Stdout, err.Error()+"\n")
			}
		},
	}

	return cmd
}
