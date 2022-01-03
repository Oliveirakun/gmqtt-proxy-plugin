package proxy

import (
	"context"
	"fmt"

	"github.com/DrmagicE/gmqtt/server"
)

func (p *Proxy) HookWrapper() server.HookWrapper {
	return server.HookWrapper{
		OnMsgArrivedWrapper: p.OnMsgArrivedWrapper,
	}
}

func (p *Proxy) OnMsgArrivedWrapper(pre server.OnMsgArrived) server.OnMsgArrived {
	return func(ctx context.Context, client server.Client, req *server.MsgArrivedRequest) error {

		message := fmt.Sprintf("Message received - Topic:  %v, Payload: %v\n", req.Message.Topic, string(req.Message.Payload))
		p.logger.Info(message)

		return nil
	}
}
