package proxy

import (
	"time"

	"github.com/DrmagicE/gmqtt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Consumer struct {
	channel   chan *gmqtt.Message
	semaphore *Semaphore
	client    mqtt.Client
}

func NewConsumer(c chan *gmqtt.Message, s *Semaphore, cl mqtt.Client) *Consumer {
	return &Consumer{
		channel:   c,
		semaphore: s,
		client:    cl,
	}
}

func (c *Consumer) SetClient(cl mqtt.Client) {
	c.client = cl
}

func (c *Consumer) DequeueMessagesAndSendToBroker() {
	for {
		if c.semaphore.IsOpen() {
			message := c.dequeue()
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

	t := c.client.Publish(message.Topic, message.QoS, false, message.Payload)
	go func() {
		<-t.Done()
		if t.Error() != nil {
			log.Error(t.Error().Error())
		}
	}()
}
