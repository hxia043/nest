package image

import (
	"github.com/hxia043/nest/internal/engine"
	"github.com/hxia043/nest/internal/nestctl/db"
)

func NewEngineOption(namespace string, registry *db.Registry) engine.Options {
	options := make(engine.Options)

	options["username"] = registry.Username
	options["password"] = registry.Password
	options["ca-cert"] = registry.CaCert
	options["registry"] = registry.Registry
	options["namespace"] = namespace

	return options
}
