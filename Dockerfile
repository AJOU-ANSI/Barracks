FROM golang:latest

RUN go get -u github.com/golang/dep/cmd/dep

RUN mkdir -p /go/src/Barracks

WORKDIR /go/src/Barracks

COPY . .

RUN dep ensure

RUN go build main.go

CMD ['./main', '-contest', 'shake16open']