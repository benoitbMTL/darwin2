#!/bin/bash

# Function to update from Git
update_from_git() {
    echo "Fetching updates from Git..."
    if ! git fetch; then
        echo "Failed to fetch updates from Git."
        exit 1
    fi

    echo "Displaying branch status..."
    git branch -v

    echo "Merging changes from origin/main..."
    if ! git merge origin/main; then
        echo "Failed to merge changes from origin/main."
        exit 1
    fi
}

# Function to serve the application (Development)
serve_app() {
    # Start Vue.js Development Server
    echo "Starting Vue.js development server..."
    cd vue
    npm run serve &
    cd ..

    # Start Go Server
    echo "Building and running Go server..."
    cd go
    if ! go run .; then
        echo "Go server failed to start."
        exit 1
    fi
}

# Function to build and run the application (Production)
build_app() {
    echo "Building Vue.js application..."
    cd vue
    if ! npm run build; then
        echo "Vue.js build failed."
        exit 1
    fi
    cd ..

    echo "Building and running Go server..."
    cd go
    if ! go build -o server; then
        echo "Go build failed."
        exit 1
    fi
    ./server
}

# Main script
case "$1" in
    -serve)
        update_from_git
        serve_app
        ;;
    -build)
        update_from_git
        build_app
        ;;
    *)
        echo "Usage: $0 -serve | -build"
        exit 1
        ;;
esac
