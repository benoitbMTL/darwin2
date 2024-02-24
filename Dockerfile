# Define the Go build stage
FROM golang:latest as go-builder
WORKDIR /go/src/app
COPY go/ .
# Copying the Go modules manifests and fetching the dependencies
COPY go/go.mod go/go.sum ./
RUN go mod download
# Building the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -v -o darwin2

# Define the Node build stage for Vue.js
FROM node:latest as vue-builder
WORKDIR /app
COPY vue/package.json vue/package-lock.json ./
# Install Vue CLI and project dependencies
RUN npm install -g @vue/cli && npm install
COPY vue/ .
# Build the Vue.js distribution
RUN npm run build

# Define the final image
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
# Copy the Vue.js distribution from the Vue builder stage
COPY --from=vue-builder /app/dist ./vue/dist
# Copy the Go binary and other necessary files from the Go builder stage
COPY --from=go-builder /go/src/app/darwin2 .

# Expose the port the app runs on
EXPOSE 8080
# Run the Go binary
CMD ["./darwin2"]
