FROM golang:1.10-alpine
#RUN apk add --no-cache curl bash git openssh

#RUN apk add --no-cache gcc musl-dev
RUN apk add --update --no-cache alpine-sdk bash python
WORKDIR /root
RUN git clone https://github.com/edenhill/librdkafka.git
WORKDIR /root/librdkafka
RUN /root/librdkafka/configure
RUN make
RUN make install
#For golang applications
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh