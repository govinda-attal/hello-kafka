FROM gattal/goalpine-librdkafka:1.10 AS build-env

ARG X_LDFLAGS_ARGS


WORKDIR /go/src/github.com/gattal/hello-kafka
COPY . .

RUN mkdir dist && \
	mkdir dist/producer

RUN dep ensure -v

RUN sleep 5

RUN	GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.version=1.1.0'" -o ./dist/producer/producer ./cmd/producer/...

FROM alpine:3.7
RUN apk -U add ca-certificates

WORKDIR /app
COPY --from=build-env /go/src/github.com/gattal/hello-kafka/dist/producer/ /app/
CMD /app/producer