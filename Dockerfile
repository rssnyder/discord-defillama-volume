FROM golang:1.19-alpine
LABEL org.opencontainers.image.source https://github.com/rssnyder/discord-defillama-volume

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /discord-bot

ENTRYPOINT /discord-bot -token "$TOKEN" -protocol "$PROTOCOL" -activity "$ACTIVITY" -refresh "${REFRESH:-30}" -metrics "${METRICS:-:8080}" -nickname
