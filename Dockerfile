# Builder Image
FROM golang:alpine3.11 AS compiler

COPY . esw-hcl-compiler
WORKDIR esw-hcl-compiler
RUN mkdir /binaries && go mod vendor && GOOS=linux go build -o esw-c && mv esw-c /binaries/

FROM alpine:3.11

ARG ESW_COMP_VERSION

LABEL maintainer="Kevin Haller <keivn.haller@tuwien.ac.at>"
LABEL version="${ESW_COMP_VERSION}"
LABEL description="Image for compiling a HCL configuration for the Exploratory Search web app."

COPY --from=compiler /binaries/esw-c /usr/local/bin/esw-c
RUN chmod a+x /usr/local/bin/esw-c

ENTRYPOINT ["esw-c"]