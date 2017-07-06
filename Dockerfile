FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN apk update && \
    apk add git curl && \
    go-wrapper download && \
    go-wrapper install

CMD ["go-wrapper", "run"]

EXPOSE 8080

HEALTHCHECK --interval=5s --timeout=3s \
  CMD curl -f http://localhost:8080 || exit 1
