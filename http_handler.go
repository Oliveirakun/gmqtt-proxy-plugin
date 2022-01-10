package proxy

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
)

type HTTPHandler struct {
	proxy *Proxy
}

func NewHTTPHandler(p *Proxy) *HTTPHandler {
	return &HTTPHandler{proxy: p}
}

func (h *HTTPHandler) Handle(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	pattern := runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1}, []string{"v1", "toogle"}, "", runtime.AssumeColonVerbOpt(true)))
	mux.Handle("PUT", pattern, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		stop := req.URL.Query().Get("stop")
		switch stop {
		case "true":
			h.proxy.closeSemaphore()
		case "false":
			brokerURI := req.URL.Query().Get("broker-uri")
			h.proxy.changeRemoteBrokerURI(brokerURI)
			h.proxy.openSemaphore()
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": stop})
	})

	return nil
}
