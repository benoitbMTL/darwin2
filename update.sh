#!/bin/bash

# Function to update from Git
update_from_git() {
    echo "Fetching updates from Git..."
    git fetch || { echo "Failed to fetch updates."; exit 1; }

    echo "Displaying branch status..."
    git branch -v

    echo "Checking for changes in the vue directory..."
    if git diff --name-only origin/main | grep -q "vue/"; then
        echo "Changes detected in the vue directory."
        vue_changes=1
    else
        echo "No changes detected in the vue directory."
        vue_changes=0
    fi

    echo "Merging changes from origin/main..."
    git merge origin/main || { echo "Failed to merge changes."; exit 1; }

    return $vue_changes
}

# Function to build Vue.js application
build_vue_app() {
    echo "Building Vue.js application..."
    cd vue || exit 1
    npm install || { echo "Vue.js installation failed."; exit 1; }
    npm run build || { echo "Vue.js build failed."; exit 1; }
    cd ..
}

# Function to build and run Go application
serve_go_app() {
    echo "Building and running Go server..."
    cd go || exit 1
    go build . || { echo "Go build failed."; exit 1; }
    ./darwin2 &
    sleep 2 # Wait for the server to start
}

# Function to manage Docker
manage_docker() {
    container_name="darwin2"

    # Check if container is running
    running_container=$(docker ps -q -f name=^/${container_name}$)
    if [[ -n "$running_container" ]]; then
        echo "Stopping and removing existing Docker container..."
        docker stop "$container_name"
        docker rm "$container_name"
    fi

    echo "Building Docker image and running container..."
    docker build -t "$container_name" .
    docker run -dp 8080:8080 --name "$container_name" "$container_name"
}

# Function to handle script options
# Added additional option: build
build_and_serve() {
    if update_from_git; then
        build_vue_app
    fi
    serve_go_app
}

# Main script execution
case $1 in
    run)
        build_and_serve
        curl -X GET http://localhost:8080/reset
        ;;
    docker)
        update_from_git
        manage_docker
        ;;
    update)
        update_from_git
        ;;
    *)
        echo "Usage: $0 {run|docker|update}"
        exit 1
        ;;
esac
