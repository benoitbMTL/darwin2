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
    ##############################################################
    echo "Fetching updates from Git..."
    git fetch || { echo "Failed to fetch updates."; exit 1; }

    echo "Displaying branch status..."
    git branch -v

    echo "Checking for changes in the vue directory..."
    if git diff --name-only origin/main | grep "vue/" | grep -v "vue/package-lock.json" | grep -v "vue/package.json" | grep -q .; then
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
    ##############################################################
    echo "Building Vue.js application..."
    cd vue || exit 1
    npm install || { echo "Vue.js installation failed."; exit 1; }
    npm run build || { echo "Vue.js build failed."; exit 1; }
    cd ..
}

# Function to build and run Go application
serve_go_app() {
    ##############################################################
    echo "Building and running Go server..."
    cd go || exit 1
    go build . || { echo "Go build failed."; exit 1; }
    ./darwin2
}

# Function to manage Docker
manage_docker() {

    # Get host mapping configuration from the first script argument or default to 5
    local config_option="${1:-5}"
    set_host_mappings "$config_option"

    # Install Docker if not present
    if ! command -v docker &> /dev/null; then
        ##############################################################
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
        ##############################################################
        echo "Stopping existing Docker container..."
        docker stop "$container_name"
        echo "Removing existing Docker container..."
        docker rm "$container_name"
    fi

    # Build a fresh Docker image
    ##############################################################
    echo "Building Docker image and running container..."
    docker build -t "$container_name" .

    # Construct the Docker run command with dynamic --add-host options
    DOCKER_RUN_CMD="docker run -dp 8080:8080 --name \"$container_name\""
    for HOST_MAPPING in "${HOST_MAPPINGS[@]}"
    do
        DOCKER_RUN_CMD+=" $HOST_MAPPING"
    done
    DOCKER_RUN_CMD+=" $container_name"
    
    # Run the Docker container with the specified host entries
    eval "$DOCKER_RUN_CMD"
}

# Function to set host mappings based on selected configuration
set_host_mappings() {
    case $1 in
        1) # FabricLab FortiWeb
            HOST_MAPPINGS=(
                "--add-host dvwa.corp.fabriclab.ca:10.163.7.23"
                "--add-host bank.corp.fabriclab.ca:10.163.7.23"
                "--add-host juiceshop.corp.fabriclab.ca:10.163.7.24"
                "--add-host petstore3.corp.fabriclab.ca:10.163.7.25"
                "--add-host speedtest.corp.fabriclab.ca:10.163.7.26"
            )
            ;;
        2) # FabricLab FortiADC
            HOST_MAPPINGS=(
                "--add-host dvwa.corp.fabriclab.ca:10.163.7.31"
                "--add-host bank.corp.fabriclab.ca:10.163.7.32"
                "--add-host juiceshop.corp.fabriclab.ca:10.163.7.33"
                "--add-host petstore3.corp.fabriclab.ca:10.163.7.34"
                "--add-host speedtest.corp.fabriclab.ca:10.163.7.35"
            )
            ;;
        3) # FabricLab FortiWeb03
            HOST_MAPPINGS=(
                "--add-host dvwa.corp.fabriclab.ca:10.163.7.41"
                "--add-host bank.corp.fabriclab.ca:10.163.7.41"
                "--add-host juiceshop.corp.fabriclab.ca:10.163.7.42"
                "--add-host petstore3.corp.fabriclab.ca:10.163.7.43"
                "--add-host speedtest.corp.fabriclab.ca:10.163.7.44"
            )
            ;;
        4) # NUC FortiWeb
            HOST_MAPPINGS=(
                "--add-host dvwa.corp.fabriclab.ca:192.168.4.10"
                "--add-host bank.corp.fabriclab.ca:192.168.4.10"
                "--add-host juiceshop.corp.fabriclab.ca:192.168.4.20"
                "--add-host petstore3.corp.fabriclab.ca:192.168.4.30"
                "--add-host speedtest.corp.fabriclab.ca:192.168.4.40"
            )
            ;;
        5|*) # No host or default
            HOST_MAPPINGS=()
            ;;
    esac
}

# Function to install required packages and setup the environment
install_environment() {
    # Install Linux packages
    ##############################################################
    print_header "Updating package lists..."
    sudo apt update || { echo "Failed to update package lists."; exit 1; }
    print_completion "Done!"

    # Install Linux packages
    ##############################################################
    print_header "Installing required Linux packages..."
    sudo apt install nmap unzip tree net-tools vim perl libnet-ssleay-perl libio-socket-ssl-perl -y || { echo "Failed to install required Linux packages."; exit 1; }
    print_completion "Done!"

    # Install Nikto
    ##############################################################
    print_header "Installing Nikto..."
    cd "$SCRIPT_DIR/go" || exit 1
    rm -Rf nikto
    git clone https://github.com/sullo/nikto.git
    cd "$SCRIPT_DIR"
    print_completion "Done!"

    # Install Go
    ##############################################################
    print_header "Installing Go..."
    echo "Downloading Go from ${GOURL}..."
    curl -L -s -O ${GOURL} || { echo "Failed to download Go."; exit 1; }
    sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf $GOPACKAGE || { echo "Failed to install Go."; exit 1; }
    rm -f ${GOPACKAGE}
    echo "Done!"
    
    # Define GO PATH
    echo "Exporting Go binary path..."
    export PATH=$PATH:/usr/local/go/bin

    # Get the current user's shell
    CURRENT_SHELL=$(echo $SHELL | awk -F'/' '{print $NF}')

    FILE_TO_UPDATE=""

    # Determine which configuration file to update based on the shell
    case $CURRENT_SHELL in
        bash)
            FILE_TO_UPDATE="$HOME/.bashrc"
            ;;
        zsh)
            FILE_TO_UPDATE="$HOME/.zshrc"
            ;;
        *)
            echo "Unsupported shell or unable to determine."
            exit 1
            ;;
    esac

    # Update the corresponding configuration file
    if [ -n "$FILE_TO_UPDATE" ]; then
        # Set GOROOT
        if ! grep -q "export GOROOT=/usr/local/go" "$FILE_TO_UPDATE"; then
            echo "export GOROOT=/usr/local/go" >> "$FILE_TO_UPDATE"
        fi

        # Set GOPATH
        if ! grep -q "export GOPATH=\$HOME/go" "$FILE_TO_UPDATE"; then
            echo "export GOPATH=\$HOME/go" >> "$FILE_TO_UPDATE"
        fi

        # Set PATH for Go
        if ! grep -q "\$GOROOT/bin:\$GOPATH/bin" "$FILE_TO_UPDATE"; then
            echo 'export PATH=$PATH:$GOROOT/bin:$GOPATH/bin' >> "$FILE_TO_UPDATE"
        fi

        echo "Configuration updated in $FILE_TO_UPDATE"
    else
        echo "No configuration file to update."
    fi

    print_completion "Done!"

    print_completion "You might need to source $FILE_TO_UPDATE or log out and back in for this change to take effect."

    # Install Node.js and npm (https://deb.nodesource.com)
    ##############################################################
    print_header "Installing Node.js..."
    sudo apt-get update && sudo apt-get install -y ca-certificates curl gnupg
    sudo rm -f /etc/apt/keyrings/nodesource.gpg
    curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_$NODE_MAJOR.x nodistro main" | sudo tee /etc/apt/sources.list.d/nodesource.list
    sudo apt-get update && sudo apt-get install nodejs -y
    print_completion "Done!"

    # Install Bootstrap and Bootstrap Icons locally within the Vue project
    ##############################################################
    print_header "Installing Bootstrap and Bootstrap Icons locally..."
    cd "$SCRIPT_DIR/vue" || exit 1
    npm install bootstrap bootstrap-icons || { echo "Failed to install Bootstrap and Bootstrap Icons."; exit 1; }
    print_completion "Done!"

    # Setup Vue.js application
    ##############################################################
    print_header "Setting up Vue.js application..."
    npm install @vue/cli || { echo "Failed to install Vue CLI."; exit 1; }
    npm install || { echo "Failed to install Vue.js application dependencies."; exit 1; }
    echo "Vue.js application setup completed. You can now run 'npm run serve' to start the application."
    cd "$SCRIPT_DIR"
    print_completion "Done!"

    # Install Chrome & Chromedriver
    ##############################################################
    # https://googlechromelabs.github.io/chrome-for-testing/

    print_header "Installing Chrome & Chromedriver..."

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

    print_completion "Done!"

    ##############################################################
    print_versions

}

# Function to print section headers in bold and color
print_header() {
    # ANSI color codes
    CYAN='\033[0;36m'
    NC='\033[0m' # No Color
    echo -e "${CYAN}\n===========================================================================${NC}"
    echo -e "${CYAN}$1${NC}"
    echo -e "${CYAN}===========================================================================${NC}\n"
}

# Function to print completion messages in bold and color
print_completion() {
    # ANSI color codes
    GREEN='\033[0;32m'
    NC='\033[0m' # No Color
    echo -e "${GREEN}---------------------------------------------------------------------------${NC}"
    echo -e "${GREEN}$1${NC}"
    echo -e "${GREEN}---------------------------------------------------------------------------${NC}\n"
}

# Function to print installed package versions
print_versions() {
    # Check for Go version
    if command -v go &> /dev/null; then
        go_version=$(go version)
    else
        go_version="Not installed"
    fi

    # Print versions
    printf "\n===========================================================================\n"
    printf "Summary of installed packages and versions:\n"
    printf "Go:\t\t\t%s\n" "$go_version"

    # Check for Nikto version
    if [ -d "$SCRIPT_DIR/go/nikto" ]; then
        nikto_version=$(perl go/nikto/program/nikto.pl -Version)
    else
        nikto_version="Not installed"
    fi

    printf "Nikto:\t\t\t%s\n" "$nikto_version"

    # Check for Node.js version
    if command -v node &> /dev/null; then
        node_version=$(node -v)
    else
        node_version="Not installed"
    fi

    printf "Node.js:\t\t%s\n" "$node_version"

    # Check for npm version
    if command -v npm &> /dev/null; then
        npm_version=$(npm -v)
    else
        npm_version="Not installed"
    fi

    printf "npm:\t\t\t%s\n" "$npm_version"

    # Check for Bootstrap version

    cd "$SCRIPT_DIR/vue"

    if npm list bootstrap &> /dev/null; then
        bootstrap_version=$(npm list bootstrap | grep bootstrap | head -1 | awk '{print $2}')
    else
        bootstrap_version="Not installed"
    fi

    printf "Bootstrap:\t\t%s\n" "$bootstrap_version"
    
    # Check for Bootstrap Icons version
    if npm list bootstrap-icons &> /dev/null; then
        bootstrap_icons_version=$(npm list bootstrap-icons | grep bootstrap-icons | head -1 | awk '{print $2}')
    else
        bootstrap_icons_version="Not installed"
    fi

    printf "Bootstrap Icons:\t%s\n" "$bootstrap_icons_version"
    
    cd "$SCRIPT_DIR"

    # Check for Google Chrome version
    if [ -f "./go/selenium/chrome-linux64/chrome" ]; then
        CHROME_VERSION=$(./go/selenium/chrome-linux64/chrome --version)
    else
        CHROME_VERSION="Not installed"
    fi
    
    printf "Google Chrome:\t\t%s\n" "$CHROME_VERSION"
    
    # Check for ChromeDriver version
    if [ -f "./go/selenium/chromedriver-linux64/chromedriver" ]; then
        CHROMEDRIVER_VERSION=$(./go/selenium/chromedriver-linux64/chromedriver --version)
    else
        CHROMEDRIVER_VERSION="Not installed"
    fi

    printf "ChromeDriver:\t\t%s\n" "$CHROMEDRIVER_VERSION"
    echo -e "===========================================================================\n"

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
    printf "Usage: %s {run|docker|update|force|install|version|help} [options]\n" "$(basename "$0")"
    printf "  %-20s%s\n" "run:" "Build and serve the application."
    printf "  %-20s%s\n" "docker:" "Manage Docker container for the application. Optional [config] to set host mappings: 1 for FabricLab FortiWeb, 2 for FabricLab FortiADC, 3 for FabricLab FortiWeb03, 4 for NUC FortiWeb, 5 for no hosts."
    printf "  %-20s%s\n" "update:" "Update the application from Git."
    printf "  %-20s%s\n" "force:" "Force build and serve the application (ignores changes)."
    printf "  %-20s%s\n" "install:" "Install and initialize environment to run the application."
    printf "  %-20s%s\n" "version:" "Print the versions of the installed packages."
    printf "  %-20s%s\n" "help:" "Display this help message."
    printf "\nExamples:\n"
    printf "  %-20s%s\n" "$(basename "$0") docker 1" "Use FabricLab FortiWeb host settings."
    printf "  %-20s%s\n" "$(basename "$0") docker 5" "Run Docker without specific host settings."
}

# Main script execution
case $1 in
    run)
        build_and_serve
        ;;
    docker)
        update_from_git
        manage_docker $2  # Pass the second command line argument to manage_docker
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
