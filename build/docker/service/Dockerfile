FROM golang:1.22.3-alpine3.18 AS build
LABEL authors="IAV"
WORKDIR /MinioFileService
COPY go.mod go.sum ./
RUN go mod download
COPY . .
WORKDIR /MinioFileService/cmd/app
RUN go build -o build

FROM alpine:3.18 AS prod
COPY --from=build /MinioFileService /.
EXPOSE 8080
ENTRYPOINT ["/cmd/app/build"]
