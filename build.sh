#!/bin/bash

echo "Fetching updates from Git..."
if ! git fetch; then
    echo "Failed to fetch updates from Git."
    exit 1
fi

echo "Displaying branch status..."
git branch -v

echo "Merging changes from origin/main..."
if ! git merge origin/main; then
    echo "Failed to merge changes from origin/main."
    exit 1
fi

echo "Building Vue.js application..."
cd vue
if ! npm run build; then
    echo "Vue.js build failed."
    exit 1
fi
cd ..

echo "Building and running Go server..."
cd go
if ! go build -o server; then
    echo "Go build failed."
    exit 1
fi
./server