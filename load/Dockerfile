FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN apk update && \
    apk add git && \
    go-wrapper download && \
    go-wrapper install && \
    mkdir -p /results

CMD ["go-wrapper", "run"]
