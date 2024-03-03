#!/bin/bash

# Save the current working directory
SCRIPT_DIR=$(pwd)

# Function to update from Git
update_from_git() {
    echo -e "\n--------------------------------------------------"
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

    echo $vue_changes
}

# Function to build Vue.js application
build_vue_app() {
    echo -e "\n--------------------------------------------------"
    echo "Building Vue.js application..."
    cd vue || exit 1
    npm install || { echo "Vue.js installation failed."; exit 1; }
    npm run build || { echo "Vue.js build failed."; exit 1; }
    cd ..
}

# Function to build and run Go application
serve_go_app() {
    echo -e "\n--------------------------------------------------"
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
        echo -e "\n--------------------------------------------------"
        echo "Stopping existing Docker container..."
        docker stop "$container_name"
        echo "Removing existing Docker container..."
        docker rm "$container_name"
    fi

    echo -e "\n--------------------------------------------------"
    echo "Building Docker image and running container..."
    docker build -t "$container_name" .
    docker run -dp 8080:8080 --name "$container_name" "$container_name"
}

# Function to install required packages and setup the environment
init_environment() {
    echo -e "\n--------------------------------------------------"
    echo "Initializing environment for the application..."

    # Install Linux packages
    echo -e "\n--------------------------------------------------"
    echo "Updating package lists..."
    sudo apt update || { echo "Failed to update package lists."; exit 1; }

    echo -e "\n--------------------------------------------------"
    echo "Installing required Linux packages..."
    sudo apt install nmap tree net-tools vim perl -y || { echo "Failed to install required Linux packages."; exit 1; }

    # Install Nikto
    echo -e "\n--------------------------------------------------"
    echo "Installing Nikto..."
    cd "$SCRIPT_DIR/go" || exit 1
    if [ ! -d "nikto" ]; then
        git clone https://github.com/sullo/nikto.git
        echo "Nikto cloned successfully."
    else
        nikto_version=$(perl nikto/program/nikto.pl -Version | grep -o 'Nikto [0-9]*\.[0-9]*\.[0-9]*')
        echo "Nikto version found: $nikto_version"
        if [ "$nikto_version" != "Nikto 2.5.0" ]; then
            echo "Warning: Nikto version is not 2.5.0. Current version: $nikto_version"
        fi
    fi
    cd "$SCRIPT_DIR"

    # Install Go
    echo -e "\n--------------------------------------------------"
    echo "Installing Go..."
    wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz || { echo "Failed to download Go."; exit 1; }
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz || { echo "Failed to install Go."; exit 1; }
    export PATH=$PATH:/usr/local/go/bin
    go_version=$(go version)
    echo "Go version: $go_version"

    # Install Node.js and npm
    # https://deb.nodesource.com/
    echo -e "\n--------------------------------------------------"
    echo "Installing Node.js..."
    sudo apt-get update && sudo apt-get install -y ca-certificates curl gnupg
    sudo rm -f /etc/apt/keyrings/nodesource.gpg
    curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
    NODE_MAJOR=20
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | sudo tee /etc/apt/sources.list.d/nodesource.list
    sudo apt-get update && sudo apt-get install nodejs -y
    node_version=$(node -v)
    echo "Node version: $node_version"

    # Install Bootstrap and Bootstrap Icons locally within the Vue project
    echo -e "\n--------------------------------------------------"
    echo "Installing Bootstrap and Bootstrap Icons locally..."
    cd "$SCRIPT_DIR/vue" || exit 1
    npm install bootstrap bootstrap-icons || { echo "Failed to install Bootstrap and Bootstrap Icons."; exit 1; }

    # Setup Vue.js application
    echo -e "\n--------------------------------------------------"
    echo "Setting up Vue.js application..."
    npm install @vue/cli || { echo "Failed to install Vue CLI."; exit 1; }
    npm install || { echo "Failed to install Vue.js application dependencies."; exit 1; }
    echo "Vue.js application setup completed. You can now run 'npm run serve' to start the application."
    cd "$SCRIPT_DIR"

    # Summarize installed package versions
    printf "\n--------------------------------------------------\n"
    printf "Summary of installed packages and versions:\n"
    printf "Go:\t\t\t%s\n" "$go_version"
    # Add similar printf statements for other software versions installed in this script
    printf "Nikto version:\t\t%s\n" "$nikto_version"
    printf "Node.js:\t\t%s\n" "$node_version"
    npm_version=$(npm -v)
    printf "npm:\t\t\t%s\n" "$npm_version"
    # Assuming Bootstrap and Bootstrap Icons versions are fetched from package.json or similar
    bootstrap_version=$(npm list bootstrap | grep bootstrap | head -1 | awk '{print $2}')
    printf "Bootstrap:\t\t%s\n" "$bootstrap_version"
    bootstrap_icons_version=$(npm list bootstrap-icons | grep bootstrap-icons | head -1 | awk '{print $2}')
    printf "Bootstrap Icons:\t%s\n" "$bootstrap_icons_version"
    printf "Environment initialization completed successfully.\n"
    printf "--------------------------------------------------\n"

}

# Function to display help
print_help() {
    printf "Usage: %s {run|docker|update|init|help}\n" "$(basename "$0")"
    printf "  %-20s%s\n" "run:" "Build and serve the application."
    printf "  %-20s%s\n" "docker:" "Manage Docker container for the application."
    printf "  %-20s%s\n" "update:" "Update the application from Git."
    printf "  %-20s%s\n" "init:" "Initialize environment to run the application."
    printf "  %-20s%s\n" "help:" "Display this help message."
}

build_and_serve() {
    vue_changes=$(update_from_git)
    if [ "$vue_changes" -eq 1 ]; then
        build_vue_app
    fi
    serve_go_app
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
