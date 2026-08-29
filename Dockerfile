ARG GOLANG_VERSION=1.26
FROM registry.access.redhat.com/ubi9/go-toolset:${GOLANG_VERSION} AS builder
WORKDIR /app
COPY maas-controller/go.mod maas-controller/go.sum ./maas-controller/
RUN cd maas-controller && go mod download
COPY maas-controller/ ./maas-controller/
RUN cd maas-controller && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -trimpath -ldflags="-s -w" -o /manager ./cmd/manager

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
WORKDIR /
COPY --from=builder /manager /manager
COPY maas-api/deploy /maas-api/deploy
COPY deployment/base/maas-api /deployment/base/maas-api
COPY deployment/base/maas-controller/policies /deployment/base/maas-controller/policies
COPY deployment/base/payload-processing /deployment/base/payload-processing
COPY deployment/components /deployment/components
RUN chmod +x /manager && chmod -R g=u /maas-api /deployment
USER 1001
ENTRYPOINT ["/manager"]
