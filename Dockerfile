FROM golang:1.18

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go env -w  GOPROXY=https://goproxy.cn,direct
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN go install github.com/swaggo/swag/cmd/swag
RUN go mod download && go mod verify

COPY . .
RUN go generate .
RUN go build -buildvcs=false -v -o /usr/local/bin/app .

CMD ["app"]