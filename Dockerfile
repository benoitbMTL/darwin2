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
RUN npm install bootstrap bootstrap-icons
# Build the Vue.js distribution
RUN npm run build

# Use Ubuntu as the base image
FROM ubuntu:latest

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

# Install necessary packages for Google Chrome, Nikto, and additional dependencies
RUN apt-get update && apt-get install -y \
    wget \
    gnupg \
    ca-certificates \
    git \
    perl \
    libnet-ssleay-perl \
    libio-socket-ssl-perl \
    nmap \
    unzip \
    tree \
    net-tools \
    vim \
    --no-install-recommends \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Set up directory for Go environment
RUN mkdir -p /go/selenium

# Download and install Chrome and Chromedriver
RUN wget https://storage.googleapis.com/chrome-for-testing-public/122.0.6261.128/linux64/chrome-linux64.zip -O /tmp/chrome.zip \
    && wget https://storage.googleapis.com/chrome-for-testing-public/122.0.6261.128/linux64/chromedriver-linux64.zip -O /tmp/chromedriver.zip \
    && unzip /tmp/chrome.zip -d /go/selenium/ \
    && unzip /tmp/chromedriver.zip -d /go/selenium/ \
    && rm /tmp/chrome.zip /tmp/chromedriver.zip

# Make Chromedriver executable
RUN chmod +x /go/selenium/chromedriver-linux64/chromedriver

# Set the environment variable for the container
ENV PATH="$PATH:/go/selenium/chrome-linux64:/go/selenium/chromedriver-linux64"

# Clone Nikto from GitHub into the specified Go directory
RUN git clone https://github.com/sullo/nikto.git /go/nikto

# Copy the Go binary to the /go directory
COPY --from=go-builder /go/src/app/darwin2 /go/darwin2

# Copy the built Vue.js application to the /vue/dist directory
COPY --from=vue-builder /app/dist /vue/dist

# Expose the port the app runs on
EXPOSE 8080

# Run the print_versions script when the container starts, followed by your application
CMD ["/bin/bash", "-c", "cd /go && ./darwin2"]
