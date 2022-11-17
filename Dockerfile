# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app
COPY . .

RUN go get -d -v
RUN go build -v

CMD ["./spreewill-core"]