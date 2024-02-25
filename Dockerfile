# Stage 1: Build Go binary
FROM golang:1.19-alpine AS go

WORKDIR /go/src/darwin2
COPY go .
RUN go mod download
RUN go build -o main

# Stage 2: Serve frontend with Nginx
FROM nginx:1.23-alpine

WORKDIR /usr/share/nginx/html
COPY ../vue/dist .
RUN echo "server_names_hash_bucket_size 64;" >> /etc/nginx/conf.d/default.conf

# Copy Go binary from previous stage
COPY --from=go /go/src/darwin2/main /usr/bin/main

CMD ["/usr/sbin/nginx", "-g", "daemon off;"]
