FROM golang:alpine as base-builder

LABEL maintainer='@ctrose17 <>'

WORKDIR /app

EXPOSE 8081

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build

FROM alpine:latest

WORKDIR /root/

COPY --from=base-builder /app/data/stateMatcher.json /root/
COPY --from=base-builder /app/data/regionMatcher.json /root/
COPY --from=base-builder /app/SimFBA .

ENV PORT 8081
ENV GOPATH /go
EXPOSE 8081

CMD ["./SimFBA"]