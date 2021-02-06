# Two-stage build:
#    first  FROM prepares a binary file in full environment ~780MB
#    second FROM takes only binary file ~20MB
FROM golang:1.15 AS builder

# File Author / Maintainer
MAINTAINER VELES GROUP

RUN apt-get update && apt-get install -y net-tools dnsutils  ca-certificates libproj-dev protobuf-compiler && apt-get clean -y
WORKDIR /root

ENV GO111MODULE=auto
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOROOT=/usr/local/go
ENV GOBIN=/root/go
ENV GOPATH $HOME/go
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin:$GOROOT:$GOPATH:$GOBIN
ENV CGO_CFLAGS="-g -O2"
ENV CGO_CPPFLAGS=""
ENV CGO_CXXFLAGS="-g -O2"
ENV CGO_FFLAGS="-g -O2"
ENV CGO_LDFLAGS="-g -O2"
#ENV GCCGO="gccgo"
#ENV CC="clang"
#ENV CXX="clang++"
ENV GOGCCFLAGS="-fPIC -m64 -pthread -fno-caret-diagnostics -Qunused-arguments -fmessage-length=0 -gno-record-gcc-switches -fno-common"

RUN cat /etc/*-release
RUN pwd
RUN ls -l /home/
RUN ls -l /root/

RUN mkdir /root/.ssh/

ADD .ssh/.gitconfig          /root/.gitconfig
ADD .ssh/config              /root/.ssh/config
ADD .ssh/id_git_dig_center   /root/.ssh/id_git_dig_center

RUN chmod 600 /root/.ssh/id_git_dig_center

RUN ssh-keyscan -t rsa git.c.dig.center 2>&1 >> /root/.ssh/known_hosts
RUN ssh-keyscan -t rsa git.c.dig.center:22022 2>&1 >> /root/.ssh/known_hosts

RUN ssh -vT git@git.c.dig.center

RUN set

RUN mkdir /app && mkdir /app/templates && mkdir /app/etc
ADD *.go /app/
# ADD go.* /app/
WORKDIR /app/
RUN cd /app

RUN go get -insecure -u git.c.dig.center/go/env.git
RUN go get -insecure -u git.c.dig.center/go/cache.git
RUN go get -insecure -u git.c.dig.center/go/auth.git

RUN go get -u github.com/jung-kurt/gofpdf

RUN go get -u google.golang.org/grpc

RUN apt-get update && apt-get install -y protobuf-compiler && apt-get clean -y
RUN go get -u github.com/golang/protobuf/protoc-gen-go

RUN find / -name "protoc"
RUN ls -l
RUN set
RUN go version

RUN mkdir /tml && \
    cd /tmp && \
    git clone git@git.c.dig.center:rpc-proto/bpmn.git && \
    cd /tmp/bpmn && ./go-compile.sh

# RUN go get -v -d ./...
RUN ls -l
RUN set
RUN go version
RUN go build -v -o ./web-service ./...

RUN ls -l /go/src/

RUN rm -rf /root/.ssh/

#########
# second stage to obtain a very small image
FROM alpine:latest
# File Author / Maintainer
MAINTAINER VELES GROUP

RUN mkdir /app && mkdir /app/static && mkdir /app/etc

VOLUME /app/etc
VOLUME /app/templates
VOLUME /app/storage

WORKDIR /app

COPY --from=builder /app/web-service /app/web-service
RUN chmod +x /app/web-service

# libproj-dev 
RUN apk update && \
    apk add -u ca-certificates && \
    rm -rf /var/lib/apt/lists/*

ADD ./docker/nsswitch.conf /etc/nsswitch.conf
ADD ./fonts       /app/fonts
#ADD ./web-service /app/web-service

# Run the command on container startup
EXPOSE 3000
CMD ["/app/web-service"]
