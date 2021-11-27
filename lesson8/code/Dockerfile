FROM golang:1.17.1-alpine3.14 as builder
WORKDIR /apps
COPY . .
RUN CGO_ENALBED=0 GOOS=linux go build -o lesson8


FROM alpine:latest as prod
WORKDIR /apps
COPY --from=builder /apps/lesson8 /apps
EXPOSE 9900
CMD ["./lesson8"]