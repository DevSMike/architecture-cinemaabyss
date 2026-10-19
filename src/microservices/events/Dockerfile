FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy the source code
COPY . .

# Resolve dependencies (fills go.sum, incl. kafka-go's transitive deps)
RUN go mod tidy && go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o events-service .

# Use a smaller image for the final application
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/events-service .

# Expose the application port
EXPOSE 8082

# Command to run the executable
CMD ["./events-service"]