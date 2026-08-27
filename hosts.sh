#!/bin/bash

# Function to remove existing entries
remove_existing_entries() {
    sed -i '/dvwa.corp.fabriclab.ca/d' /etc/hosts
    sed -i '/bank.corp.fabriclab.ca/d' /etc/hosts
    sed -i '/juiceshop.corp.fabriclab.ca/d' /etc/hosts
    sed -i '/petstore3.corp.fabriclab.ca/d' /etc/hosts
    sed -i '/speedtest.corp.fabriclab.ca/d' /etc/hosts
}

# Check permissions for hosts file
if [ ! -w /etc/hosts ]; then
  echo "Error: You do not have permission to modify the hosts file. Please run the script with sudo."
  exit 1
fi

echo "Choose an option:"
echo "1) fortiweb"
echo "2) fortiweb2"
echo "3) fortiadc"
read -p "Enter your choice (1, 2, or 3): " choice

# Remove existing entries to prevent duplicates
remove_existing_entries

# Add new entries based on the choice
case $choice in
    1)  # fortiweb
        echo "10.163.7.23 dvwa.corp.fabriclab.ca" >> /etc/hosts
        echo "10.163.7.23 bank.corp.fabriclab.ca" >> /etc/hosts
        echo "10.163.7.24 juiceshop.corp.fabriclab.ca" >> /etc/hosts
        echo "10.163.7.25 petstore3.corp.fabriclab.ca" >> /etc/hosts
        echo "10.163.7.26 speedtest.corp.fabriclab.ca" >> /etc/hosts
        ;;
    2)  # fortiweb2
        echo "10.163.7.41 dvwa.corp.fabriclab.ca" >> /etc/hosts
        echo "10.163.7.41 bank.corp.fabriclab.ca" >> /etc/hosts
        echo "10.163.7.42 juiceshop.corp.fabriclab.ca" >> /etc/hosts
        echo "10.163.7.43 petstore3.corp.fabriclab.ca" >> /etc/hosts
        echo "10.163.7.44 speedtest.corp.fabriclab.ca" >> /etc/hosts
        ;;
    3)  # fortiadc
        echo "10.163.7.31 dvwa.corp.fabriclab.ca" >> /etc/hosts
        echo "10.163.7.32 bank.corp.fabriclab.ca" >> /etc/hosts
        echo "10.163.7.33 juiceshop.corp.fabriclab.ca" >> /etc/hosts
        echo "10.163.7.34 petstore3.corp.fabriclab.ca" >> /etc/hosts
        echo "10.163.7.35 speedtest.corp.fabriclab.ca" >> /etc/hosts
        ;;
    *)
        echo "Invalid option selected."
        exit 2
        ;;
esac

# Verify entries were added
echo "The following entries were added to the hosts file:"
grep -e "dvwa.corp.fabriclab.ca" -e "bank.corp.fabriclab.ca" -e "juiceshop.corp.fabriclab.ca" -e "petstore3.corp.fabriclab.ca" -e "speedtest.corp.fabriclab.ca" /etc/hosts
