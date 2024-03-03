#!/bin/bash

# File to store the output
output_file="conf.txt"

# Remove the output file if it exists
rm -f "$output_file"

# List of files to process
files=(
vue/src/components/02-bot-mitigation/BotDeception.vue
vue/src/components/02-bot-mitigation/KnownBots.vue
vue/src/components/01-web-protection/CookieSecurity.vue
go/routes/routes.go
go/handlers/botDeception.go
go/handlers/botKnownBots.go
go/handlers/webCookieSecurity.go
)

# Loop through the files and append their content to conf.txt
for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo "========================================================" >> "$output_file"
        echo "File content of $file" >> "$output_file"
        echo "========================================================" >> "$output_file"
        cat "$file" >> "$output_file"
        echo "" >> "$output_file" # Add an extra newline for readability
    else
        echo "Warning: File $file not found." >> "$output_file"
    fi
done

echo "========================================================" >> "$output_file"

echo "Done!"
