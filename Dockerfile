FROM golang:alpine3.12 as build
RUN apk add build-base && apk add --no-cache ca-certificates && update-ca-certificates
WORKDIR server
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -ldflags "-linkmode external -extldflags -static" -o /server/bin/server ./cmd/ontology

FROM scratch
COPY --from=build /server/bin/server .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["./server"]
