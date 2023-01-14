package registry

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hxia043/nest/internal/swicctl/db"
	"github.com/spf13/cobra"
)

type registryUpdateOption struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Registry string `json:"registry"`
	certPath string
	Cert     string `json:"cert"`
}

func newRegistryUpdateOption() *registryUpdateOption {
	return &registryUpdateOption{}
}

func (o *registryUpdateOption) validateRegistryOption() error {
	if o.Username == "" && o.Password == "" && o.Registry == "" && o.certPath == "" {
		return fmt.Errorf("at least one username, password, registry or certpath needed for update")
	}

	return nil
}

func (o *registryUpdateOption) getUpdateRegistry(registry *db.Registry) *registryUpdateOption {
	updateRegistry := registryUpdateOption{Name: o.Name}

	if o.Username != "" {
		updateRegistry.Username = o.Username
	} else {
		updateRegistry.Username = registry.Username
	}

	if o.Password != "" {
		updateRegistry.Password = o.Password
	} else {
		updateRegistry.Password = registry.Password
	}

	if o.certPath != "" {
		updateRegistry.Cert = readCertFrom(o.certPath)
	} else {
		updateRegistry.Cert = registry.CaCert
	}

	if o.Registry != "" {
		updateRegistry.Registry = o.Registry
	} else {
		updateRegistry.Registry = registry.Registry
	}

	return &updateRegistry
}

func (o *registryUpdateOption) run(args []string) error {
	var err error
	o.Name, err = getRegistryName(args)
	if err != nil {
		return err
	}

	if err := o.validateRegistryOption(); err != nil {
		return err
	}

	dbOption, err := db.InitDB()
	if err != nil {
		return err
	}

	registry, err := dbOption.Get(o.Name)
	if err != nil {
		return err
	}

	updateRegistry := o.getUpdateRegistry(registry)
	data, err := json.Marshal(updateRegistry)
	if err != nil {
		return err
	}

	if err := dbOption.Update(data); err != nil {
		return err
	}

	return nil
}

func NewUpdateRegistryCmd() *cobra.Command {
	o := newRegistryUpdateOption()
	cmd := &cobra.Command{
		Use:   "registry [name]",
		Short: "Update registry",
		Long:  "Update registry",
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

	cmd.Flags().StringVarP(&o.Username, "username", "u", "", "Username for login registry")
	cmd.Flags().StringVarP(&o.Password, "password", "p", "", "Password for login registry")
	cmd.Flags().StringVarP(&o.Registry, "registry", "r", "", "Kubernetes registry")
	cmd.Flags().StringVarP(&o.certPath, "cert", "c", "", "Cert path for login registry")

	return cmd
}
