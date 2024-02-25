# Define the Go build stage
FROM golang:latest as go-builder
WORKDIR /build
COPY go/ .
# Fetching the dependencies
RUN go mod download
# Building the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -v -o darwin2

# Define the Node build stage for Vue.js
FROM node:latest as vue-builder
WORKDIR /app
COPY vue/ .
# Install Vue CLI and project dependencies
RUN npm install -g @vue/cli && npm install
RUN npm install bootstrap
RUN npm install bootstrap-icons
# Build the Vue.js distribution
RUN npm run build

# Define the final image
FROM alpine:latest  
WORKDIR /go
COPY --from=go-builder /build/darwin2 .
WORKDIR /vue
COPY --from=vue-builder /app/dist .

# Change back to the go directory to run the application
WORKDIR /go

# Expose the port the app runs on
EXPOSE 8080

# Run the Go binary
CMD ["./darwin2"]
