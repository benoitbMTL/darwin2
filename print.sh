#!/bin/bash

# Remove the conf.txt file if it exists
rm -f conf.txt

# Traverse the "./vue/src" directory for .vue files
echo "====================" >> conf.txt
find ./vue/src -type f -name "*.vue" -exec echo "{}" >> conf.txt \;
for file in $(find ./vue/src -type f -name "*.vue"); do
    echo "====================" >> conf.txt
    echo "File content of $file" >> conf.txt
    cat "$file" >> conf.txt
done

# Add contents of vue/src/router/index.js
echo "====================" >> conf.txt
echo "File content of vue/src/router/index.js" >> conf.txt
cat "vue/src/router/index.js" >> conf.txt

# Add contents of vue/src/main.js
echo "====================" >> conf.txt
echo "File content of vue/src/main.js" >> conf.txt
cat "vue/src/main.js" >> conf.txt

# Traverse the "./go" directory for .go files
echo "====================" >> conf.txt
find ./go -type f -name "*.go" -exec echo "{}" >> conf.txt \;
for file in $(find ./go -type f -name "*.go"); do
    echo "====================" >> conf.txt
    echo "File content of $file" >> conf.txt
    cat "$file" >> conf.txt
done
echo "====================" >> conf.txt

echo "Done!"
