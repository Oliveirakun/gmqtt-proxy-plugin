package proxy

import (
	"sync"
	"time"

	"github.com/DrmagicE/gmqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Consumer struct {
	channel    chan *gmqtt.Message
	semaphore  *Semaphore
	client     mqtt.Client
	clientLock sync.RWMutex
}

func NewConsumer(c chan *gmqtt.Message, s *Semaphore, cl mqtt.Client) *Consumer {
	return &Consumer{
		channel:   c,
		semaphore: s,
		client:    cl,
	}
}

func (c *Consumer) SetClient(cl mqtt.Client) {
	c.clientLock.Lock()
	defer c.clientLock.Unlock()

	c.client = cl
}

func (c *Consumer) getClient() mqtt.Client {
	c.clientLock.RLock()
	defer c.clientLock.RUnlock()

	return c.client
}

func (c *Consumer) DequeueMessagesAndSendToBroker() {
	for {
		if c.semaphore.IsOpen() {
			message := c.dequeue()
			time.Sleep(200 * time.Millisecond)
			c.sendToBroker(message)
		}
	}
}

func (c *Consumer) dequeue() *gmqtt.Message {
	select {
	case m := <-c.channel:
		return m
	case <-time.After(1 * time.Second):
		return nil
	}
}

func (c *Consumer) sendToBroker(message *gmqtt.Message) {
	if message == nil {
		return
	}

	t := c.getClient().Publish(message.Topic, message.QoS, false, message.Payload)
	go func() {
		<-t.Done()
		if t.Error() != nil {
			log.Error(t.Error().Error())
		}
	}()
}
