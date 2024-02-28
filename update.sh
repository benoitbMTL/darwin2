#!/bin/bash

# Save the current working directory
SCRIPT_DIR=$(pwd)

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
    ./darwin2
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

# Function to install required packages and setup the environment
init_environment() {
    echo "Initializing environment for the application..."

    # Install Linux packages
    echo "Updating package lists..."
    sudo apt update || { echo "Failed to update package lists."; exit 1; }

    echo "Installing required Linux packages (nmap, tree, net-tools, vim)..."
    sudo apt install nmap tree net-tools vim -y || { echo "Failed to install required Linux packages."; exit 1; }

    # Install Nikto
    echo "Installing Nikto..."
    # Navigate to the 'go' directory from 'darwin2'
    cd "$SCRIPT_DIR/go" || exit 1
    if [ ! -d "nikto" ]; then
        git clone https://github.com/sullo/nikto.git
        echo "Nikto cloned successfully."
    else
        version=$(perl nikto/program/nikto.pl -Version | grep -o 'Nikto [0-9]*\.[0-9]*\.[0-9]*')
        echo "Nikto version found: $version"
        if [ "$version" != "Nikto 2.5.0" ]; then
            echo "Warning: Nikto version is not 2.5.0. Current version: $version"
        fi
    fi
    # Return to the 'darwin2' directory
    cd "$SCRIPT_DIR" || exit 1

    # Install Go
    echo "Installing Go..."
    wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz || { echo "Failed to download Go."; exit 1; }
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz || { echo "Failed to install Go."; exit 1; }
    export PATH=$PATH:/usr/local/go/bin
    echo "Go version: $(go version)"

    # Install Node.js and npm
    echo "Installing Node.js..."
    # The Node.js installation steps already assume return to the 'darwin2' directory

    # Install Bootstrap and Bootstrap Icons globally
    echo "Installing Bootstrap and Bootstrap Icons globally..."
    npm install -g bootstrap bootstrap-icons || { echo "Failed to install Bootstrap and Bootstrap Icons."; exit 1; }

    # Setup Vue.js application
    echo "Setting up Vue.js application..."
    cd "$SCRIPT_DIR/vue" || exit 1
    npm install @vue/cli || { echo "Failed to install Vue CLI."; exit 1; }
    npm install || { echo "Failed to install Vue.js application dependencies."; exit 1; }
    echo "Vue.js application setup completed. You can now run 'npm run serve' to start the application."

    # Return to the 'darwin2' directory
    cd "$SCRIPT_DIR" || exit 1

    echo "Environment initialization completed successfully."
}

# Function to display help
print_help() {
    echo "Usage: $0 {run|docker|update|init|help}"
    echo "  run: Build and serve the application."
    echo "  docker: Manage Docker container for the application."
    echo "  update: Update the application from Git."
    echo "  init: Initialize environment to run the application."
    echo "  help: Display this help message."
}

# Main script execution
case $1 in
    run)
        build_and_serve
        ;;
    docker)
        update_from_git
        manage_docker
        ;;
    update)
        update_from_git
        ;;
    init)
        init_environment
        ;;
    help)
        print_help
        ;;
    *)
        print_help
        exit 1
        ;;
esac
