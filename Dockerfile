FROM golang:alpine AS builder
WORKDIR /app
COPY ./ ./
WORKDIR  /app
RUN ls
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/main /app/
RUN ls
RUN chmod +x /app/main
ENTRYPOINT ["/app/main"]
EXPOSE 8080
#FROM ubuntu
# MAINTAINER Ayush
#WORKDIR goapps

#COPY  ./Databases/main  /goapps/main
#RUN chmod +x /goapps/main
# ENV PORT 8080
#EXPOSE 8080
#ENTRYPOINT /goapps/main