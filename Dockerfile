ARG VERSION=latest

FROM golang:1.18 as builder
ARG VERSION

WORKDIR /go/src

COPY . /go/src

RUN go get
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -o ecs -ldflags "-s -w -X main.version=${VERSION}"
RUN md5sum ecs

FROM amazon/aws-cli:2.7.9

# Install ssm plugin for exec
RUN curl "https://s3.amazonaws.com/session-manager-downloads/plugin/latest/linux_64bit/session-manager-plugin.rpm" -o "session-manager-plugin.rpm"
RUN yum install -y session-manager-plugin.rpm

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/ecs /ecs
COPY --from=builder /tmp /tmp

ENTRYPOINT [ "/ecs" ]
