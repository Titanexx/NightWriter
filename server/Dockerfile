FROM golang:1.16-alpine as base

WORKDIR /nightwriter

COPY go.mod .
COPY go.sum .
RUN go mod download

FROM base as debug
CMD ["go","run","."]  

FROM base as builder
WORKDIR /nightwriter
COPY . .
RUN go build .

FROM alpine:latest as prod
WORKDIR /nightwriter
COPY --from=builder /nightwriter/nightwriter .
CMD [ "./nightwriter" ]

FROM alpine:latest as prod-debug
WORKDIR /nightwriter
COPY --from=builder /nightwriter/nightwriter .
CMD [ "./nightwriter", "-d"]