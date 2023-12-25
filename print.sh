#!/bin/bash

# Print the content of App.vue
echo "#######################################"
echo "Content of vue/src/App.vue:"
echo "#######################################"
cat vue/src/App.vue
echo "#######################################"

# Print the contents of all .vue files in the views folder
for file in vue/src/views/*.vue; do
    echo "Content of $file:"
    echo "#######################################"
    cat "$file"
    echo "#######################################"
done

# Print the content of index.js in the router folder
echo "Content of vue/src/router/index.js:"
echo "#######################################"
cat vue/src/router/index.js
echo "#######################################"

# Print the content of index.html in the public folder
echo "Content of vue/public/index.html:"
echo "#######################################"
cat vue/public/index.html
echo "#######################################"