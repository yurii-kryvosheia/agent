package unix

import (
	"github.com/grafana/agent/component"
	"github.com/grafana/agent/pkg/integrations"
	node_integration "github.com/grafana/agent/pkg/integrations/node_exporter"
)

func init() {
	component.Register(component.Registration{
		Name:      "prometheus.exporter.unix",
		Args:      Config{},
		Exports:   Exports{},
		Singleton: true,
		Build:     New(createExporter, "unix"),
	})
}

func createExporter(opts component.Options, args component.Arguments) (integrations.Integration, error) {
	cfg := args.(Config)
	return node_integration.New(opts.Logger, cfg.Convert())
}
