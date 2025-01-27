package cynxhostcentral

import "cynxhostagent/internal/dependencies"

type CynxhostCentral struct {
	Config *dependencies.ConfigCentral
}

func New(config *dependencies.ConfigCentral) *CynxhostCentral {
	return &CynxhostCentral{
		Config: config,
	}
}

func (c *CynxhostCentral) CallShutdownCallback() error {
	return nil
}
