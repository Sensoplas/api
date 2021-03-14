FROM golang:1.16.2-alpine AS build

WORKDIR /build
COPY . .

RUN go build -o /bin/api

FROM alpine:3.13 AS runtime

COPY --from=build /bin/api /bin/api
COPY --from=build /etc/nsswitch.conf /etc/nsswitch.conf

LABEL app="sensoplas-api"

ENV PORT="8080"

ENTRYPOINT [ "/bin/api", "http" ]