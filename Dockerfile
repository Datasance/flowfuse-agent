FROM golang:1.20.14-alpine AS go-builder

ARG TARGETOS
ARG TARGETARCH

RUN mkdir -p /go/src/github.com/datasance/flowfuse-agent
WORKDIR /go/src/github.com/datasance/flowfuse-agent
COPY . /go/src/github.com/datasance/flowfuse-agent
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o bin/flowfuse-agent


FROM node:18-alpine

ARG VERSION=latest
COPY LICENSE /licenses/LICENSE
COPY --from=go-builder /go/src/github.com/datasance/flowfuse-agent/bin/flowfuse-agent /bin/flowfuse-agent

RUN apk add --no-cache --virtual buildtools build-base linux-headers udev python3 openssl

RUN mkdir -m 777 -p /opt/flowfuse-device
RUN npm config set cache /opt/flowfuse-device/.npm --global
RUN npm install -g @flowfuse/device-agent@${VERSION}
RUN chmod -R 777 /opt/flowfuse-device/.npm

EXPOSE 1880

LABEL org.label-schema.name="FlowFuse Device Agent" \
    org.label-schema.url="https://flowfuse.com" \
    org.label-schema.vcs-type="Git" \
    org.label-schema.vcs-url="https://github.com/FlowFuse/device-agent" \
    org.label-schema.docker.dockerfile="docker/Dockerfile" \
    org.schema-label.description="Collaborative, low code integration and automation environment" \
    authors="FlowFuse Inc."

CMD ["/bin/flowfuse-agent"]
