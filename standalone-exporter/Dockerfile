FROM golang:1.27.3 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 go build -o /out/tiny-exporter .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/tiny-exporter /tiny-exporter
ENTRYPOINT ["/tiny-exporter"]
