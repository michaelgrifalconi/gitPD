# gitpd build container TODO: automate job to update base images
FROM golang:1.12.1-alpine3.9 AS build-env

RUN apk add --no-cache --upgrade git openssh-client ca-certificates
RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/michaelgrifalconi/gitpd
#COPY Gopkg.toml Gopkg.lock ./
COPY . / ./
RUN dep ensure -vendor-only -v
#COPY *.go ./
RUN go build -v -o /go/bin/gitpd ./cmd/gitpd/main.go

# Final container
FROM alpine:3.9.2
RUN apk add --no-cache --upgrade git openssh-client

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /go/bin/gitpd /usr/bin/gitpd

#ENV PATH="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" Need this?

# Create a generic SSH config for Github
WORKDIR /root/.ssh
COPY build/builder/sshconfig config

WORKDIR /root/
ENTRYPOINT [ "gitpd" ]