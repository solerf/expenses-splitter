FROM golang:1.24.3-alpine3.21 AS build

RUN go version

RUN apk update; apk add make

RUN mkdir app-build

WORKDIR /go/app-build

COPY ./accounting ./accounting
COPY ./httpx ./httpx
COPY ./go.mod ./main.go ./Makefile ./

RUN make tests
RUN make build

FROM alpine:3.21.3 AS final

RUN mkdir -p /bill-splitter
COPY --from=build /go/app-build/bin/bill-splitter /bill-splitter/app

WORKDIR /bill-splitter
EXPOSE 8000
CMD [ "./app" ]
