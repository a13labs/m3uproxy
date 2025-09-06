# Stage 1: Build the Go binary
FROM golang:1.24.0-alpine AS go-builder

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download the module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .
RUN chmod +x scripts/*

# Build the Go binaries
RUN go build -o m3uproxy server/main.go
RUN go build -o m3uproxy-cli cli/main.go

# Stage 2: Build the React SPA
FROM node:22-alpine AS spa-builder

# Set the working directory
WORKDIR /app

# Copy the React application code
COPY ./player /app

# Install dependencies and build the project
RUN npm install 
ARG NODE_ENV=production
RUN npm run build

# Create a zip file of the SPA build output
RUN apk add --no-cache zip && \
    zip -jr /app/player.zip /app/dist

# Stage 3: Final stage, build the final image
FROM alpine:latest

# Copy Go binaries and other files from the go-builder stage
COPY --from=go-builder /app/m3uproxy /app/m3uproxy
COPY --from=go-builder /app/m3uproxy-cli /app/m3uproxy-cli
COPY --from=go-builder /app/conf /app/conf
COPY --from=go-builder /app/scripts/entrypoint.sh /app/entrypoint.sh

# Copy the SPA zip file from spa-builder stage
COPY --from=spa-builder /app/player.zip /app/assets/player.zip

# Add necessary packages
RUN apk --no-cache add ca-certificates ffmpeg bash && \
    mkdir -p /app/cache 

# Set the working directory
WORKDIR /app

# Expose the application port
EXPOSE 8080

# Set the entrypoint script
CMD ["/app/entrypoint.sh"]
