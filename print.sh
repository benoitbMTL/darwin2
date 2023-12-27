#!/bin/bash

# Step 1: Remove the conf.txt file if it exists
rm -f conf.txt

# Step 2: Traverse the "./vue/src" directory for .vue files
echo "====================" >> conf.txt
find ./vue/src -type f -name "*.vue" -exec echo "{}" >> conf.txt \;
for file in $(find ./vue/src -type f -name "*.vue"); do
    echo "====================" >> conf.txt
    echo "File content of $file" >> conf.txt
    cat "$file" >> conf.txt
done
echo "====================" >> conf.txt

# Step 3: Traverse the "./go" directory for .go files
find ./go -type f -name "*.go" -exec echo "{}" >> conf.txt \;
for file in $(find ./go -type f -name "*.go"); do
    echo "====================" >> conf.txt
    echo "File content of $file" >> conf.txt
    cat "$file" >> conf.txt
done
echo "====================" >> conf.txt

echo "Done!"
