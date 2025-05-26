FROM golang:1.24.3
LABEL authors="iceesar@live.com"
WORKDIR /app
COPY . ./
RUN GOOS=linux GOARCH=amd64 go build -o my-app
ENTRYPOINT ["./my-app"]