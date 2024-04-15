#!/bin/bash

# Check permissions for hosts file
if [ ! -w /etc/hosts ]; then
  echo "Error: You do not have permission to modify the hosts file. Run the script with sudo."
  exit 1
fi

# Add entries to hosts file
echo "" >> /etc/hosts
echo "# Add entries for FabricLab sites" >> /etc/hosts
echo "dvwa.corp.fabriclab.ca:10.163.7.23" >> /etc/hosts
echo "bank.corp.fabriclab.ca:10.163.7.23" >> /etc/hosts
echo "juiceshop.corp.fabriclab.ca:10.163.7.24" >> /etc/hosts
echo "petstore3.corp.fabriclab.ca:10.163.7.25" >> /etc/hosts
echo "speedtest.corp.fabriclab.ca:10.163.7.26" >> /etc/hosts

# Verify entries were added
echo "The following entries were added to the hosts file:"
cat /etc/hosts | grep -e "dvwa.corp.fabriclab.ca" -e "bank.corp.fabriclab.ca" -e "juiceshop.corp.fabriclab.ca" -e "petstore3.corp.fabriclab.ca" -e "speedtest.corp.fabriclab.ca"
