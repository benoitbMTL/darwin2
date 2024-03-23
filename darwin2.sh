#!/bin/bash

# Initialize global variables
SCRIPT_DIR=$(pwd)
NODE_MAJOR=20
vue_changes=0
GOPACKAGE="go1.22.0.linux-amd64.tar.gz"
GOURL="https://go.dev/dl/${GOPACKAGE}"
# Download locations: version 122.0.6261.128
CHROMEDRIVER_URL="https://storage.googleapis.com/chrome-for-testing-public/122.0.6261.128/linux64/chromedriver-linux64.zip"
CHROME_URL="https://storage.googleapis.com/chrome-for-testing-public/122.0.6261.128/linux64/chrome-linux64.zip"

# Function to update from Git
update_from_git() {
    echo -e "\n---------------------------------------------------------------------------"
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
    echo -e "\n---------------------------------------------------------------------------"
    echo "Building Vue.js application..."
    cd vue || exit 1
    npm install || { echo "Vue.js installation failed."; exit 1; }
    npm run build || { echo "Vue.js build failed."; exit 1; }
    cd ..
}

# Function to build and run Go application
serve_go_app() {
    echo -e "\n---------------------------------------------------------------------------"
    echo "Building and running Go server..."
    cd go || exit 1
    go build . || { echo "Go build failed."; exit 1; }
    ./darwin2
}

# Function to manage Docker
manage_docker() {

    # Install Docker if not present
    if ! command -v docker &> /dev/null; then
        echo -e "\n---------------------------------------------------------------------------"
        echo "Docker is not installed. Installing Docker..."
        sudo apt-get update -y
        sudo apt-get upgrade -y
        sudo apt-get install apt-transport-https ca-certificates curl software-properties-common -y
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
        sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
        sudo apt-get update
        sudo apt-get install docker-ce -y
        if ! getent group docker > /dev/null; then
            sudo groupadd docker
        fi
        if [ "$(whoami)" != "root" ]; then
            sudo usermod -aG docker "$(whoami)"
            echo "Please log out and then log back in for the changes to take effect."
        fi
        echo "Docker installation and setup completed."
    fi

    # Check if a container is already running
    container_name="darwin2"
    running_container=$(docker ps -q -f name=^/${container_name}$)
    if [[ -n "$running_container" ]]; then
        echo -e "\n---------------------------------------------------------------------------"
        echo "Stopping existing Docker container..."
        docker stop "$container_name"
        echo "Removing existing Docker container..."
        docker rm "$container_name"
    fi

    # Build a fresh Docker image
    echo -e "\n---------------------------------------------------------------------------"
    echo "Building Docker image and running container..."
    docker build -t "$container_name" .
    docker run -dp 8080:8080 --name "$container_name" "$container_name"
}

# Function to install required packages and setup the environment
install_environment() {
    echo -e "\n---------------------------------------------------------------------------"
    echo "Initializing environment for the application..."

    # Install Linux packages
    echo -e "\n---------------------------------------------------------------------------"
    echo "Updating package lists..."
    sudo apt update || { echo "Failed to update package lists."; exit 1; }

    echo -e "\n---------------------------------------------------------------------------"
    echo "Installing required Linux packages..."
    sudo apt install nmap unzip tree net-tools vim perl libnet-ssleay-perl libio-socket-ssl-perl -y || { echo "Failed to install required Linux packages."; exit 1; }

    # Install Nikto
    echo -e "\n---------------------------------------------------------------------------"
    echo "Installing Nikto..."
    cd "$SCRIPT_DIR/go" || exit 1
    rm -Rf nikto
    git clone https://github.com/sullo/nikto.git
    echo "Done!"

    cd "$SCRIPT_DIR"

    # Install Go
    echo -e "\n---------------------------------------------------------------------------"
    echo "Installing Go..."
    echo "Downloading Go from ${GOURL}..."
    curl -L -s -O ${GOURL} || { echo "Failed to download Go."; exit 1; }
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf $GOPACKAGE || { echo "Failed to install Go."; exit 1; }
    rm -f ${GOPACKAGE}
    export PATH=$PATH:/usr/local/go/bin


    # Define GO PATH
    GOPATH="export PATH=\$PATH:/usr/local/go/bin"

    # Check if GOPATH already exists in .bashrc
    if grep -Fxq "$GOPATH" ~/.bashrc; then
        echo "GOPATH already exists in .bashrc."
    else
        # If GOPATH doesn't exist, add it to .bashrc
        echo "$GOPATH" >> ~/.bashrc
        echo "GOPATH has been added to .bashrc."
        echo "Please log off and log back in for the changes to take effect."
    fi

    # Install Node.js and npm
    # https://deb.nodesource.com/
    echo -e "\n---------------------------------------------------------------------------"
    echo "Installing Node.js..."
    sudo apt-get update && sudo apt-get install -y ca-certificates curl gnupg
    sudo rm -f /etc/apt/keyrings/nodesource.gpg
    curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | sudo tee /etc/apt/sources.list.d/nodesource.list
    sudo apt-get update && sudo apt-get install nodejs -y


    # Install Bootstrap and Bootstrap Icons locally within the Vue project
    echo -e "\n---------------------------------------------------------------------------"
    echo "Installing Bootstrap and Bootstrap Icons locally..."
    cd "$SCRIPT_DIR/vue" || exit 1
    npm install bootstrap bootstrap-icons || { echo "Failed to install Bootstrap and Bootstrap Icons."; exit 1; }


    # Setup Vue.js application
    echo -e "\n---------------------------------------------------------------------------"
    echo "Setting up Vue.js application..."
    npm install @vue/cli || { echo "Failed to install Vue CLI."; exit 1; }
    npm install || { echo "Failed to install Vue.js application dependencies."; exit 1; }
    echo "Vue.js application setup completed. You can now run 'npm run serve' to start the application."
    cd "$SCRIPT_DIR"

    # Install Chrome & Chromedriver
    # https://googlechromelabs.github.io/chrome-for-testing/

    echo -e "\n---------------------------------------------------------------------------"
    echo "Installing Chrome & Chromedriver..."

    cd vue
    npm install selenium-webdriver
    cd "$SCRIPT_DIR"

    # Create a directory for the setup
    cd go
    rm -Rf ./selenium
    mkdir -p ./selenium
    cd selenium

    # Download ChromeDriver and Chrome
    wget $CHROMEDRIVER_URL -O chromedriver.zip
    wget $CHROME_URL -O chrome.zip

    # Unzip the downloaded files
    unzip chromedriver.zip
    unzip chrome.zip

    # Clean up the setup directory
    rm -rf *.zip
    cd "$SCRIPT_DIR"

    echo "Done!"

    print_versions

}

# Function to print installed package versions
print_versions() {
    # Check for Go version
    if command -v go &> /dev/null; then
        go_version=$(go version)
    else
        go_version="Not installed"
    fi

    # Check for Nikto version
    if [ -d "$SCRIPT_DIR/go/nikto" ]; then
        nikto_version=$(perl go/nikto/program/nikto.pl -Version | grep -o 'Nikto [0-9]*\.[0-9]*\.[0-9]*')
    else
        nikto_version="Not installed"
    fi

    # Check for Node.js version
    if command -v node &> /dev/null; then
        node_version=$(node -v)
    else
        node_version="Not installed"
    fi

    # Check for npm version
    if command -v npm &> /dev/null; then
        npm_version=$(npm -v)
    else
        npm_version="Not installed"
    fi

    # Check for Bootstrap version

    cd "$SCRIPT_DIR/vue"

    if npm list bootstrap &> /dev/null; then
        bootstrap_version=$(npm list bootstrap | grep bootstrap | head -1 | awk '{print $2}')
    else
        bootstrap_version="Not installed"
    fi

    # Check for Bootstrap Icons version
    if npm list bootstrap-icons &> /dev/null; then
        bootstrap_icons_version=$(npm list bootstrap-icons | grep bootstrap-icons | head -1 | awk '{print $2}')
    else
        bootstrap_icons_version="Not installed"
    fi
    
    cd "$SCRIPT_DIR"

    # Check for ChromeDriver version
    if [ -f "./go/selenium/chromedriver-linux64/chromedriver" ]; then
        CHROMEDRIVER_VERSION=$(./go/selenium/chromedriver-linux64/chromedriver --version)
    else
        CHROMEDRIVER_VERSION="Not installed"
    fi

    # Check for Google Chrome version
    if [ -f "./go/selenium/chrome-linux64/chrome" ]; then
        CHROME_VERSION=$(./go/selenium/chrome-linux64/chrome --version)
    else
        CHROME_VERSION="Not installed"
    fi

    # Print versions
    printf "\n---------------------------------------------------------------------------\n"
    printf "Summary of installed packages and versions:\n"
    printf "Go:\t\t\t%s\n" "$go_version"
    printf "Nikto:\t\t\t%s\n" "$nikto_version"
    printf "Node.js:\t\t%s\n" "$node_version"
    printf "npm:\t\t\t%s\n" "$npm_version"
    printf "Bootstrap:\t\t%s\n" "$bootstrap_version"
    printf "Bootstrap Icons:\t%s\n" "$bootstrap_icons_version"
    printf "Google Chrome:\t\t%s\n" "$CHROME_VERSION"
    printf "ChromeDriver:\t\t%s\n" "$CHROMEDRIVER_VERSION"
    printf "Environment checks completed.\n"
    echo -e "---------------------------------------------------------------------------\n"

    if [ "$nikto_version" != "Nikto 2.5.0" ]; then
        echo "Warning: Nikto version is not 2.5.0. Current version: $nikto_version"
    fi
}

# Build and serve function
build_and_serve() {
    update_from_git
    if [ "$vue_changes" -eq 1 ]; then
        build_vue_app
    fi
    serve_go_app
}

# Force build and serve function
force_build_and_serve() {
    update_from_git
    build_vue_app
    serve_go_app
}

# Function to display help
print_help() {
    printf "Usage: %s {run|docker|update|force|install|version|help}\n" "$(basename "$0")"
    printf "  %-20s%s\n" "run:" "Build and serve the application."
    printf "  %-20s%s\n" "docker:" "Manage Docker container for the application."
    printf "  %-20s%s\n" "update:" "Update the application from Git."
    printf "  %-20s%s\n" "force:" "Force build and serve the application (ignores changes)."
    printf "  %-20s%s\n" "install:" "Install and initialize environment to run the application."
    printf "  %-20s%s\n" "version:" "Print the versions of the installed packages."
    printf "  %-20s%s\n" "help:" "Display this help message."
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
    force)
        force_build_and_serve
        ;;
    install)
        install_environment
        build_vue_app
        serve_go_app
        ;;
    version)
        print_versions
        ;;
    help)
        print_help
        ;;
    *)
        print_help
        exit 1
        ;;
esac
