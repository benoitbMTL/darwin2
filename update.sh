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

    echo "Checking for changes in the vue directory..."
    # Check if there are updates to be merged from origin/main
    if git diff --name-only origin/main | grep -q "vue/"; then
        echo "Changes detected in the vue directory."
        local vue_changes=1
    else
        echo "No changes detected in the vue directory."
        local vue_changes=0
    fi

    echo "Merging changes from origin/main..."
    if ! git merge origin/main; then
        echo "Failed to merge changes from origin/main."
        exit 1
    fi

    return $vue_changes
}

# Function to build Vue.js application and serve Go application
build_and_serve_app() {
    # Building Vue.js Application
    # Check if there were changes in the vue directory
    if [ "$(update_from_git)" -eq 1 ]; then
        echo "Building Vue.js application due to changes..."
        cd vue
        if ! npm run build; then
            echo "Vue.js build failed."
            exit 1
        fi
        cd ..
    else
        echo "Skipping Vue.js build; no changes detected."
    fi

    # Running Go Server
    echo "Running Go server..."
    cd go
    go run .
}

# Main script execution
build_and_serve_app
