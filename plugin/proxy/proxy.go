package proxy

import (
	"go.uber.org/zap"

	"github.com/DrmagicE/gmqtt/config"
	"github.com/DrmagicE/gmqtt/server"
)

var _ server.Plugin = (*Proxy)(nil)

const Name = "proxy"

func init() {
	server.RegisterPlugin(Name, New)
	config.RegisterDefaultPluginConfig(Name, &DefaultConfig)
}

func New(config config.Config) (server.Plugin, error) {
	return &Proxy{}, nil
}

var log *zap.Logger

type Proxy struct {
	logger *zap.Logger
}

func (p *Proxy) Load(service server.Server) error {
	log = server.LoggerWithField(zap.String("plugin", Name))
	p.logger = log

	apiRegistrar := service.APIRegistrar()
	err := apiRegistrar.RegisterHTTPHandler(new(HTTPHandler).Handle)
	if err != nil {
		log.Error("Error registering HTTP handler")
	}

	return nil
}

func (p *Proxy) Unload() error {
	return nil
}

func (p *Proxy) Name() string {
	return Name
}
