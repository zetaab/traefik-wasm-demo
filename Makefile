build:
	@tinygo build -o header.wasm -scheduler=none --no-debug -target=wasi header.go
