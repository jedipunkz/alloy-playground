# Build Alloy from upstream source with prometheus.exporter.tiny added.
# Build context is this directory (alloy-fork/).
FROM golang:1.27.3 AS build
ARG ALLOY_VERSION=v1.21.1
RUN git clone --depth 1 --branch ${ALLOY_VERSION} https://github.com/grafana/alloy.git /src/alloy
WORKDIR /src/alloy
COPY tiny internal/component/prometheus/exporter/tiny
COPY all_tiny.go internal/component/all/tiny.go
# Same entrypoint as upstream `make alloy`, without the UI build (SKIP_UI_BUILD) and cgo.
# gore2regex is always added by the upstream Makefile.
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    cd collector && CGO_ENABLED=0 go build -tags "netgo gore2regex" -o /out/alloy .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/alloy /bin/alloy
ENTRYPOINT ["/bin/alloy"]
