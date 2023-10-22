# Wasm middleware for Traefik

This plugin is based on [http-wasm](https://github.com/http-wasm/http-wasm-guest-tinygo)

Traefik modification: https://github.com/traefik/traefik/compare/v3.0...mmatur:traefik:feat/wasm

## Build the plugin

```bash
make build
```

## Use with Traefik

```bash
# In the new terminal
git clone git@github.com:mmatur/traefik.git
cd traefik/
git checkout feat/wasm

# Copy wasm middleware file
cp ../header.wasm .

# Create static configuration
cat <<EOF > static.yaml
entryPoints:
  web:
    address: ':8000'
  traefik:
    address: ':9080'

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
EOF

# Create dynamic configuration 
cat <<EOF > dynamic.yaml
http:
  routers:
    customer1:
      rule: Host(\`powpow.demo.traefiklabs.tech\`)
      service: customer1
      middlewares:
        - testWasm

  services:
    customer1:
      loadbalancer:
        servers:
          - url: "http://127.0.0.1:8081"

  middlewares:
    testWasm:
      wasm:
        path: ./header.wasm
EOF

# Start a whoami container
docker run -tid -p0.0.0.0:8081:80 traefik/whoami

# Run traefik
go run ./cmd/traefik/ --configFile=static.yaml
```

```bash
$ curl -H "Host: powpow.demo.traefiklabs.tech" 127.0.0.1:8000
Hostname: b8af906d1e0f
IP: 127.0.0.1
IP: 172.17.0.2
RemoteAddr: 192.168.65.1:30443
GET / HTTP/1.1
Host: powpow.demo.traefiklabs.tech
User-Agent: curl/8.1.2
Accept: */*
Accept-Encoding: gzip
X-Forwarded-For: 127.0.0.1
X-Forwarded-Host: powpow.demo.traefiklabs.tech
X-Forwarded-Port: 80
X-Forwarded-Proto: http
X-Forwarded-Server: Michaels-MacBook-Pro.local
X-Powpow: hello
X-Real-Ip: 127.0.0.1
```
# traefik-wasm-demo
