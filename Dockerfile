# -----------------------------------------------------------------------------
# Stages
# -----------------------------------------------------------------------------

ARG IMAGE_GO_BUILDER=golang:1.20.0
ARG IMAGE_FINAL=senzing/senzingapi-runtime:3.4.2

# -----------------------------------------------------------------------------
# Stage: go_builder
# -----------------------------------------------------------------------------

# define where we need to copy senzing files from
FROM ${IMAGE_FINAL} as senzing-runtime
FROM ${IMAGE_GO_BUILDER} as go_builder
ENV REFRESHED_AT=2023-02-22
LABEL Name="senzing/demo-entity-search-builder" \
      Maintainer="support@senzing.com" \
      Version="0.0.5"

# Build arguments.

ARG PROGRAM_NAME="unknown"
ARG BUILD_VERSION=0.0.0
ARG BUILD_ITERATION=0
ARG GO_PACKAGE_NAME="unknown"

# Copy local files from the Git repository.

COPY ./rootfs /
COPY . ${GOPATH}/src/${GO_PACKAGE_NAME}

# Copy necessary Senzing files from DockerHub.

COPY --from=senzing-runtime  "/opt/senzing/g2/lib/"   "/opt/senzing/g2/lib/"
COPY --from=senzing-runtime  "/opt/senzing/g2/sdk/c/" "/opt/senzing/g2/sdk/c/"

# Set path to Senzing libs.

ENV LD_LIBRARY_PATH=/opt/senzing/g2/lib/

# Build go program.

WORKDIR ${GOPATH}/src/${GO_PACKAGE_NAME}
RUN make build

# --- Test go program ---------------------------------------------------------

# Run unit tests.

# RUN go get github.com/jstemmer/go-junit-report \
#  && mkdir -p /output/go-junit-report \
#  && go test -v ${GO_PACKAGE_NAME}/... | go-junit-report > /output/go-junit-report/test-report.xml

# Copy binaries to /output.

RUN mkdir -p /output \
      && cp -R ${GOPATH}/src/${GO_PACKAGE_NAME}/target/*  /output/

# -----------------------------------------------------------------------------
# Stage: final
# -----------------------------------------------------------------------------

FROM ${IMAGE_FINAL} as final
ENV REFRESHED_AT=2023-02-22
LABEL Name="senzing/demo-entity-search" \
      Maintainer="support@senzing.com" \
      Version="0.0.5"

# Copy files from prior step.

COPY --from=go_builder "/output/linux/demo-entity-search" "/app/demo-entity-search"

# Runtime environment variables.

ENV LD_LIBRARY_PATH=/opt/senzing/g2/lib/

# Runtime execution.

WORKDIR /app
ENTRYPOINT ["/app/demo-entity-search"]
