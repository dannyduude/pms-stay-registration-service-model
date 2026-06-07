FROM harbor.ext.hp.com/pib-sw/golang:1.26.3 AS artifact-build

WORKDIR /go-quickstart

ARG ghe_auth_token
RUN if [ -z ${ghe_auth_token} ]; then \
		echo "ERROR: You need pass a 'ghe_auth_token' as argument of 'docker build' command"; \
		exit 1; \
	fi && \
	git config --global url."https://${ghe_auth_token}@github.azc.ext.hp.com/".insteadOf "https://github.azc.ext.hp.com/" && \
	go env -w GOPRIVATE=github.azc.ext.hp.com/3DSoftware/*

# Speedup build using mod download
COPY ./go.mod .
COPY ./go.sum .

COPY . .

RUN git submodule update --init --recursive

RUN go mod download

RUN CGO_ENABLED=0 go build -ldflags "-w -s" -trimpath -o /go/bin/go-quickstart_app ./app/server

# --------------------------------------------------------------
# Get a statically-linked busybox
# --------------------------------------------------------------
FROM harbor.ext.hp.com/pib-sw/alpine:3.23.0 AS busybox-source
RUN apk add --no-cache busybox-static

# --------------------------------------------------------------
# Build deployment image
# --------------------------------------------------------------

# This results in a single layer image
FROM gcr.io/distroless/static-debian11


# Copy the statically-linked busybox (works without a dynamic linker)
COPY --from=busybox-source /bin/busybox.static /bin/busybox


# See http://label-schema.org/rc1/
LABEL org.label-schema.vendo="HP Inc."
LABEL org.label-schema.name ="orchestration-service"
LABEL org.label-schema.url.description="Orchestration service"

ARG SEMANTIC_VERSION
ENV QST_SERVICE_VERSION=${SEMANTIC_VERSION}

COPY ./config /usr/src/go-quickstart/config
COPY --from=artifact-build /go/bin/go-quickstart_app /usr/src/go-quickstart

WORKDIR /usr/src/go-quickstart

ENTRYPOINT ["./go-quickstart_app", "server"]
