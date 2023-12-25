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

# Start Vue.js Development Server
echo "Starting Vue.js development server..."
cd vue
npm run serve &
cd ..

# Start Go Server
echo "Building and running Go server..."
cd go
if ! go run .; then
    echo "Go server failed to start."
    exit 1
fi
