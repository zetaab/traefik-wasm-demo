package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/http-wasm/http-wasm-host-go/handler"
	wasm "github.com/http-wasm/http-wasm-host-go/handler/nethttp"
	"github.com/rs/zerolog"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World"))
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

	b, err := json.Marshal(map[string]string{"X-foo": "Hello, World!"})
	if err != nil {
		return nil, err
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	opts := []handler.Option{
		handler.Logger(initWasmLogger(&logger)),
		handler.GuestConfig(b),
	}

	mw, err := wasm.NewMiddleware(context.Background(), code, opts...)
	if err != nil {
		return nil, err
	}
	return mw.NewHandler(context.Background(), next), nil
}
