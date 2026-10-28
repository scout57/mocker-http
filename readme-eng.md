# scout55/mocker-http

A lightweight HTTP mock server in Go. It reads YAML templates and responds to incoming requests with stubs — handy for local development and testing.

## Features

- Responses from YAML templates in `./templates`
- Exact path matching and wildcard `*` (one segment)
- Exact match takes priority over wildcard
- `response.status` (defaults to `200`)
- Response body: string, object, or array (objects/arrays are serialized to JSON)
- `Content-Type: application/json` by default
- No caching — templates are re-read on every request
- Fallback `hello world` if no matches are found
- Standard library only + `gopkg.in/yaml.v3`

## Template format

All `./templates/*.yaml` and `./templates/*.yml` files are merged. Each file is an array of objects:

```yaml
- method: GET              # case-insensitive
  path: /users/*           # * — one segment
  response:
    status: 201            # optional, default 200
    headers:               # optional
      X-Foo: bar
    body: ...              # string | object | array
```

## Usage via Docker-compose
```yaml
  scout_mocker_http:
    image: scout55/mocker
    volumes:
      - ./templates:/go/app/templates:ro
    ports:
      - 8080:8080
```

## Usage via Docker
```bash
docker run -p 8080:8080 \
  -v "$(pwd)/templates:/app/templates:ro" \
  --rm scout55/mocker
```

## Local development

```bash
# Start the environment
docker compose up -d

# Dependencies
docker compose exec app sh -c ". install.sh"

# Hot-reload
docker compose exec app air

# Shut down
docker compose down
```