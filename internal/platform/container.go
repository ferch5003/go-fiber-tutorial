package platform

import (
	"github.com/ferch5003/go-fiber-tutorial/config"
)

type Container interface {
	// CreateOrUseContainer creates or use a docker container.
	CreateOrUseContainer(config *config.EnvVars) (err error)

	// CleanContainer removes a docker container.
	CleanContainer() (err error)
}
