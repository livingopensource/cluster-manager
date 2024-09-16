FROM golang:1.22.7-alpine as build
RUN apk --update add build-base

WORKDIR /src/app
COPY . .
RUN go mod tidy

RUN CGO_ENABLED=0 go build -o /go/bin/constellation

FROM gcr.io/distroless/static-debian11
ENV TZ=Africa/Lusaka
WORKDIR /src/app

COPY --from=build /go/bin/constellation .
CMD ["/src/app/constellation"]