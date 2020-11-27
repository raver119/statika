FROM golang:1-buster as builder

# copy source files
COPY ./service /service

# get dependencies
RUN cd /service && go get -v -t -d ./...

# build app
RUN cd sources && go build -v .


FROM ubuntu:20.04

COPY --from=builder /service/service /application/statika

EXPOSE 8080
CMD cd /application && ./statika