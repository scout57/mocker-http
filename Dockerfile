# ---------- build ----------
FROM golang:1.27-alpine AS build

WORKDIR /src

COPY ./app/go.mod ./app/go.sum* ./
RUN go mod download

COPY ./app .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app .

# ---------- runtime ----------
FROM alpine:3.20

RUN adduser -D -H -u 10001 app

WORKDIR /app

COPY --from=build /out/app /app/app
COPY --from=build /src/templates /app/templates

USER app

ENTRYPOINT ["/app/app"]