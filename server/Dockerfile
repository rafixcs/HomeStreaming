FROM golang:alpine3.21 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY ./src ./src

RUN CGO_ENABLED=0 GOOS=linux go build -v -o homestreaming ./src/cmd

RUN apk add --no-cache ca-certificates

################# Debug/Dev Image #######################

FROM golang:alpine3.21 AS debug

WORKDIR /app

EXPOSE 8080

COPY go.mod go.sum ./

RUN go mod download && go mod verify

RUN go install github.com/go-delve/delve/cmd/dlv@latest

COPY ./src ./src
RUN go build -v -o /usr/local/bin/app ./src/cmd

### Run the Delve debugger ###
COPY ./dlv.sh /
RUN chmod +x /dlv.sh
ENTRYPOINT [ "/dlv.sh" ]

################# Production Image #######################

FROM golang AS prd

WORKDIR /app

EXPOSE 8080

COPY --from=build /app/homestreaming /app/

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

RUN apt-get update && apt-get install -y ffmpeg

ENTRYPOINT [ "/app/homestreaming" ]
