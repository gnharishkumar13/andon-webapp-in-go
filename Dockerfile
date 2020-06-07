
#build stage
FROM golang:1.14-alpine AS builder
WORKDIR /app/src
EXPOSE 3000
