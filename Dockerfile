FROM golang:1.16-buster as builder

# copy source files
COPY ./ /service

WORKDIR /service

# get dependencies
RUN go get -v -t -d ./...

# build app
RUN go build -v .

FROM ubuntu:20.04

# some certificates to be present
RUN apt update && apt install -y ca-certificates

COPY --from=builder /service/statika /application/statika

# setup user that will be running the server
RUN groupadd -r user && useradd -r -g user user
RUN chown -R user.user /application

# setup the default storage volume
RUN mkdir /statika && chown user:user /statika
VOLUME /statika

USER user
WORKDIR /application
CMD ["./statika"]