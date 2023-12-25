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

# Function to build Vue.js application and serve Go application
build_and_serve_app() {
    # Building Vue.js Application
    echo "Building Vue.js application..."
    cd vue
    if ! npm run build; then
        echo "Vue.js build failed."
        exit 1
    fi
    cd ..

    # Building and Running Go Server
    echo "Building and running Go server..."
    cd go
    if ! go run .; then
        echo "Go server failed to start."
        exit 1
    fi
}

# Main script execution
update_from_git
build_and_serve_app
