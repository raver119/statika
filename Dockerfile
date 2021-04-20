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

RUN groupadd -r user && useradd -r -g user user
RUN chown -R user.user /application
USER user

EXPOSE 8080
CMD cd /application && ./statika