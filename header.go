package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/http-wasm/http-wasm-guest-tinygo/handler"
	"github.com/http-wasm/http-wasm-guest-tinygo/handler/api"
)

// Config the plugin configuration.
type Config struct {
	Headers map[string]string `json:"headers,omitempty"`
}

type Middleware struct {
	Config *Config
	Start  time.Time
}

var mw = &Middleware{}

func init() {
	err := json.Unmarshal(handler.Host.GetConfig(), &mw.Config)
	if err != nil {
		handler.Host.Log(api.LogLevelError, fmt.Sprintf("Could not load config %v", err))
		os.Exit(1)
	}
}

func main() {
	handler.HandleRequestFn = mw.handleRequest
	handler.HandleResponseFn = mw.handleResponse
}

// handleRequest implements a simple request middleware.
func (mw *Middleware) handleRequest(req api.Request, resp api.Response) (next bool, reqCtx uint32) {
	handler.Host.Log(api.LogLevelDebug, "hello from handleRequest debug")
	mw.Start = time.Now()
	handler.Host.Log(api.LogLevelDebug, "time is "+mw.Start.String())
	for k, v := range mw.Config.Headers {
		req.Headers().Add(k, v)
	}
	// proceed to the next handler on the host.
	next = true
	return
}

// handleResponse implements a simple response middleware.
func (mw *Middleware) handleResponse(_ uint32, req api.Request, resp api.Response, _ bool) {
	handler.Host.Log(api.LogLevelDebug, time.Since(mw.Start).String())
}
