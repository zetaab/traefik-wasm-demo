# Wasm middleware for Traefik

This plugin is based on [http-wasm](https://github.com/http-wasm/http-wasm-guest-tinygo)

Traefik modification: https://github.com/traefik/traefik/compare/v3.0...zetaab:traefik:feature/httpwasm

## Build the plugin

```bash
make build
```

## Use with Traefik

```bash
# In the new terminal
git clone https://github.com/zetaab/traefik.git
cd traefik/
git checkout feature/httpwasm

# Create static configuration
cat <<EOF > static.yaml
entryPoints:
  web:
    address: :8000

log:
  level: debug

api:
  dashboard: true
  insecure: true

providers:
  file:
    filename: ./dynamic.yaml

metrics:
  prometheus: {}

experimental:
  localPlugins:
    wasmLocalExample:
      moduleName: github.com/zetaab/traefik-wasm-demo
EOF

# Create dynamic configuration 
cat <<EOF > dynamic.yaml
http:
  routers:
    customer1:
      rule: Host(`powpow.demo.traefiklabs.tech`)
      service: customer1
      middlewares:
        - localWasm
      
  services:
    customer1:
      loadbalancer:
        servers:
          - url: "http://127.0.0.1:8081"

  middlewares:
    localWasm:
      plugin:
        wasmLocalExample:
          headers:
            local: foo
EOF

# Start a whoami container
docker run -tid -p0.0.0.0:8081:80 traefik/whoami

# Run traefik
go run ./cmd/traefik/ --configFile=static.yaml
```

```bash
% curl -H "Host: powpow.demo.traefiklabs.tech" 127.0.0.1:8000
Hostname: 28cf73dc0e09
IP: 127.0.0.1
IP: 172.17.0.3
RemoteAddr: 172.17.0.1:48202
GET / HTTP/1.1
Host: powpow.demo.traefiklabs.tech
User-Agent: curl/8.1.2
Accept: */*
Accept-Encoding: gzip
Foo: Bar
Heh: hah
Huh: foobar
X-Forwarded-For: 127.0.0.1
X-Forwarded-Host: powpow.demo.traefiklabs.tech
X-Forwarded-Port: 80
X-Forwarded-Proto: http
X-Forwarded-Server: xx
X-Real-Ip: 127.0.0.1
```
