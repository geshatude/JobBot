FROM golang:1.26rc2-bookworm

WORKDIR /app

COPY go.mod /app
COPY go.sum /app
RUN go mod download
COPY cmd/ /app/cmd/
COPY internal/ /app/internal/

RUN CGO_ENABLED=0 GOOS=linux go build -o /jobbot /app/cmd/main.go

CMD ["/jobbot"]