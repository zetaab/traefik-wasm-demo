package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/http-wasm/http-wasm-host-go/handler"
	wasm "github.com/http-wasm/http-wasm-host-go/handler/nethttp"
	"github.com/rs/zerolog"
	"github.com/tetratelabs/wazero"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	for name, values := range r.Header {
		// Loop over all values for the name.
		for _, value := range values {
			w.Write([]byte(fmt.Sprintf("%s: %s\n", name, value)))
		}
	}
	time.Sleep(1 * time.Second)
	w.Write([]byte("Hello, World\n"))
}

func main() {
	handler, err := makeWasmHandler(http.HandlerFunc(helloWorld))
	if err != nil {
		panic(err)
	}
	http.Handle("/hello", handler)
	http.ListenAndServe(":8090", nil)
}

func makeWasmHandler(next http.Handler) (http.Handler, error) {
	code, err := os.ReadFile("../header.wasm")
	if err != nil {
		return nil, err
	}

	config := map[string]interface{}{
		"headers": map[string]string{
			"X-foo": "Hello, World!",
		},
	}

	b, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	opts := []handler.Option{
		handler.ModuleConfig(wazero.NewModuleConfig().WithSysWalltime()),
		handler.Logger(initWasmLogger(&logger)),
		handler.GuestConfig(b),
	}

	mw, err := wasm.NewMiddleware(context.Background(), code, opts...)
	if err != nil {
		return nil, err
	}
	return mw.NewHandler(context.Background(), next), nil
}
