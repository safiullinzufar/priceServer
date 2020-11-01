FROM golang:alpine

RUN apk update && apk upgrade && \
    apk add --no-cache bash git

WORKDIR /home/priceServer
COPY . .

RUN go build
CMD ["./priceServer"]
