# Docker image for GOS CLI
FROM alpine:3.18

# Install ca-certificates and other dependencies
RUN apk --no-cache add ca-certificates curl git

WORKDIR /root/

# Copy the binary
COPY gos .

# Make it executable and create symbolic link
RUN chmod +x ./gos && \
    ln -s /root/gos /usr/local/bin/gos

# Set the entrypoint
ENTRYPOINT ["gos"]

# Default command
CMD ["--help"]

# Metadata
LABEL org.opencontainers.image.title="GOS CLI"
LABEL org.opencontainers.image.description="A comprehensive Go version manager CLI"
LABEL org.opencontainers.image.vendor="Cristobal Contreras"
LABEL org.opencontainers.image.source="https://github.com/cristobalcontreras/homebrew-gos"
LABEL org.opencontainers.image.documentation="https://github.com/cristobalcontreras/homebrew-gos/blob/main/README.md"
