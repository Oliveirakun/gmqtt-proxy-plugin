package proxy

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/DrmagicE/gmqtt"
	"github.com/DrmagicE/gmqtt/config"
	"github.com/DrmagicE/gmqtt/server"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var _ server.Plugin = (*Proxy)(nil)

const Name = "proxy"

func init() {
	server.RegisterPlugin(Name, New)
	config.RegisterDefaultPluginConfig(Name, &DefaultConfig)
}

func New(config config.Config) (server.Plugin, error) {
	messageChannel := make(chan *gmqtt.Message, 10000)
	p := NewProducer(messageChannel)
	s := NewSemaphore()
	c := NewConsumer(messageChannel, s, createMQTTClient())
	go c.DequeueMessagesAndSendToBroker()

	proxy := &Proxy{
		producer:  p,
		consumer:  c,
		semaphore: s,
	}
	return proxy, nil
}

var log *zap.Logger

type Proxy struct {
	producer  *Producer
	consumer  *Consumer
	semaphore *Semaphore
}

func (p *Proxy) Load(service server.Server) error {
	log = server.LoggerWithField(zap.String("plugin", Name))

	apiRegistrar := service.APIRegistrar()
	handler := NewHTTPHandler(p.semaphore)
	err := apiRegistrar.RegisterHTTPHandler(handler.Handle)
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

func createMQTTClient() mqtt.Client {
	brokerURI := os.Getenv("REMOTE_MQTT_BROKER")
	client := mqtt.NewClient(mqtt.NewClientOptions().AddBroker(brokerURI))

	if err := client.Connect().Error(); err != nil {
		log.Error(fmt.Sprintf("Failed connecting to broker: %v", err))
	}

	return client
}
