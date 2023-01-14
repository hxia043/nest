package registry

import (
	"fmt"
	"os"

	"github.com/gosuri/uitable"
	"github.com/hxia043/nest/internal/nestctl/db"
	"github.com/hxia043/nest/pkg/util/cmdtable"
	"github.com/spf13/cobra"
)

type registryGetOption struct{}

func newRegistryGetOption() *registryGetOption {
	return &registryGetOption{}
}

func (o *registryGetOption) run(args []string) error {
	name, err := getRegistryName(args)
	if err != nil {
		return err
	}

	dbOption, err := db.InitDB()
	if err != nil {
		return err
	}

	registry, err := dbOption.Get(name)
	if err != nil {
		return fmt.Errorf("no registry %s found", name)
	}

	table := uitable.New()
	table.AddRow("NAME", "USERNAME", "REGISTRY")
	table.AddRow(registry.Name, registry.Username, registry.Registry)

	if err := cmdtable.EncodeTable(os.Stdout, table); err != nil {
		return err
	}

	return nil
}

func NewGetRegistryCmd() *cobra.Command {
	o := newRegistryGetOption()
	cmd := &cobra.Command{
		Use:     "registry [name]",
		Short:   "Get registry",
		Long:    "Get registry",
		Example: "nestctl get registry test",
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
