package proxy

import "github.com/DrmagicE/gmqtt"

type Producer struct {
	channel chan *gmqtt.Message
}

func NewProducer(c chan *gmqtt.Message) *Producer {
	return &Producer{
		channel: c,
	}
}

func (p *Producer) Enqueue(message *gmqtt.Message) {
	p.channel <- message
}
