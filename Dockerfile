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

# Use Ubuntu as the base image
FROM ubuntu:latest

# Avoid warnings by switching to noninteractive
ENV DEBIAN_FRONTEND=noninteractive

# Install necessary packages for Google Chrome and Nikto, including Perl SSL support
RUN apt-get update && apt-get install -y \
    wget \
    gnupg \
    ca-certificates \
    perl \
    git \
    libnet-ssleay-perl \
    libcrypt-ssleay-perl \
    libio-socket-ssl-perl \
    cpanminus \
    build-essential \  # Includes 'make' and other compilation tools
    libssl-dev \       # OpenSSL development packages required by Net::SSLeay
    --no-install-recommends \
    && cpanm Net::SSLeay IO::Socket::SSL \
    # Download and install Google Chrome
    && wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb \
    && apt install -y ./google-chrome-stable_current_amd64.deb \
    && rm ./google-chrome-stable_current_amd64.deb \
    # Clean up
    && apt-get purge --auto-remove -y wget gnupg cpanminus build-essential \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Verify Google Chrome installation
RUN google-chrome --version

# Copy the Go binary to the /go directory
COPY --from=go-builder /go/src/app/darwin2 /go/darwin2

# Copy the built Vue.js application to the /vue/dist directory
COPY --from=vue-builder /app/dist /vue/dist

# Copy chromedriver from your repository to the expected directory in the Docker image
COPY go/selenium/chromedriver /selenium/chromedriver

# Make chromedriver executable
RUN chmod +x /selenium/chromedriver

# Set the environment variable for the container
ENV CHROMEDRIVER_PATH="/selenium/chromedriver"

# Clone Nikto from GitHub
RUN git clone https://github.com/sullo/nikto.git /nikto
RUN perl /nikto/program/nikto.pl -Version

# Expose the port the app runs on
EXPOSE 8080

# Switch back to dialog for any ad-hoc use of apt-get
ENV DEBIAN_FRONTEND=dialog

# Run the Go binary
CMD ["/go/darwin2"]
