FROM golang:1.16-buster as builder

# copy source files
COPY ./ /service

WORKDIR /service

# get dependencies
RUN go get -v -t -d ./...

# build app
RUN go build -v .

FROM ubuntu:20.04

COPY --from=builder /service/statika /application/statika

RUN groupadd -r user && useradd -r -g user user
RUN chown -R user.user /application
USER user

EXPOSE 8080
CMD cd /application && ./statika