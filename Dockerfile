# Define the Go build stage
FROM golang:latest as go-builder
WORKDIR /go/src/app
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
# Install Perl and Git
RUN apk --no-cache add perl git

# Create the directories for Go and Vue.js applications
RUN mkdir /go && mkdir /vue
# Copy the Go binary to the /go directory
COPY --from=go-builder /go/src/app/darwin2 /go/darwin2
# Copy the built Vue.js application to the /vue/dist directory
COPY --from=vue-builder /app/dist /vue/dist

# Clone Nikto from GitHub
RUN git clone https://github.com/sullo/nikto.git /nikto

# Expose the port the app runs on
EXPOSE 8080

# Run the Go binary
CMD ["/go/darwin2"]
