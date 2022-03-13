package proxy

import (
	"fmt"
	"os"
	"time"

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
	handler := NewHTTPHandler(p)
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

func (p *Proxy) changeRemoteBrokerURI(uri string) {
	if len(uri) == 0 {
		return
	}

	newClient := createMQTTClientWithURI(uri)
	p.consumer.SetClient(newClient)
}

func (p *Proxy) openSemaphore() {
	p.semaphore.Open()
}

func (p *Proxy) closeSemaphore() {
	p.semaphore.Close()
}

func createMQTTClient() mqtt.Client {
	brokerURI := os.Getenv("REMOTE_MQTT_BROKER")
	return createMQTTClientWithURI(brokerURI)
}

func createMQTTClientWithURI(uri string) mqtt.Client {
	options := mqtt.NewClientOptions()
	options.ConnectRetry = true
	options.ConnectTimeout = 10 * time.Second
	options.ConnectRetryInterval = 2 * time.Second

	client := mqtt.NewClient(options.AddBroker(uri))

	t := client.Connect()
	if err := t.Error(); err != nil {
		log.Error(fmt.Sprintf("Failed connecting to broker: %v", err))
	}

	t.WaitTimeout(10 * time.Second)
	return client
}
