package registry

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hxia043/nest/internal/swicctl/db"
	"github.com/spf13/cobra"
)

type registryAddOption struct {
	Name       string `json:"name"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Registry   string `json:"registry"`
	caCertPath string
	CaCert     string `json:"ca-cert"`
}

func newRegistryAddOption() *registryAddOption {
	return &registryAddOption{}
}

func (o *registryAddOption) validateRegistryOption() error {
	if o.Username == "" || o.Password == "" {
		return fmt.Errorf("no username or password found")
	}

	return nil
}

func (o *registryAddOption) run(args []string) error {
	var err error
	o.Name, err = getRegistryName(args)
	if err != nil {
		return err
	}

	if err = o.validateRegistryOption(); err != nil {
		return err
	}

	dbOption, err := db.InitDB()
	if err != nil {
		return err
	}

	o.CaCert = readCertFrom(o.caCertPath)
	data, err := json.Marshal(o)
	if err != nil {
		return err
	}

	if err := dbOption.Add(data); err != nil {
		return err
	}

	return nil
}

func NewAddRegistryCmd() *cobra.Command {
	o := newRegistryAddOption()
	cmd := &cobra.Command{
		Use:     "registry [name]",
		Short:   "Add registry",
		Long:    "Add registry",
		Example: "swicctl add registry test -u nokiaadmin -p Nokia@1865 -r image-registry.openshift-image-registry.svc:5000 -c /var/home/core/cert/ca.crt",
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
	cmd.Flags().StringVarP(&o.Registry, "registry", "r", "image-registry.openshift-image-registry.svc:5000", "Kubernetes registry")
	cmd.Flags().StringVarP(&o.caCertPath, "ca-cert", "c", "/etc/docker/certs.d/image-registry.openshift-image-registry.svc:5000/ca.crt", "Cert path for login registry")

	return cmd
}
