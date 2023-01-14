package registry

import (
	"fmt"
	"os"

	"github.com/gosuri/uitable"
	"github.com/hxia043/nest/internal/nestctl/db"
	"github.com/hxia043/nest/pkg/util/cmdtable"
	"github.com/spf13/cobra"
)

type registryShowOption struct{}

func newRegistryShowOption() *registryShowOption {
	return &registryShowOption{}
}

func (o *registryShowOption) run(args []string) error {
	dbOption, err := db.InitDB()
	if err != nil {
		return err
	}

	registries, err := dbOption.Show()
	if err != nil {
		return err
	}

	table := uitable.New()
	table.AddRow("NAME", "USERNAME", "REGISTRY")
	for _, registry := range registries {
		table.AddRow(registry.Name, registry.Username, registry.Registry)
	}

	if err := cmdtable.EncodeTable(os.Stdout, table); err != nil {
		return err
	}

	return nil
}

func NewShowRegistryCmd() *cobra.Command {
	o := newRegistryShowOption()
	cmd := &cobra.Command{
		Use:   "registry",
		Short: "Show registries",
		Long:  "Show registries",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 0 {
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
