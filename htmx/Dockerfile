# Auto-generated Dockerfile file.
# See https://gowebly.org for more information.

FROM golang:1.22-alpine AS builder

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy your code into the container.
COPY . .

# Set necessary environment variables and build your project.
ENV CGO_ENABLED=0 GIN_MODE=release
RUN go build -ldflags="-s -w" -o comigo

FROM scratch

# Copy project's binary and templ from /build to the scratch container.
COPY --from=builder /build/comigo /
COPY --from=builder /build/static /static


# Set entry point.
ENTRYPOINT ["/comigo"]
