FROM gattal/goalpine-librdkafka:1.10 AS build-env

ARG X_LDFLAGS_ARGS


WORKDIR /go/src/github.com/govinda-attal/hello-kafka
COPY . .

RUN mkdir dist && \
	mkdir dist/consumer

RUN dep ensure -v

RUN	GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.version=1.1.0'" -o ./dist/consumer/consumer ./cmd/consumer/...

FROM gattal/alpine-librdkafka:0.11.5
RUN apk -U add ca-certificates

WORKDIR /app
COPY --from=build-env /go/src/github.com/govinda-attal/hello-kafka/dist/consumer/ /app/
CMD /app/consumer