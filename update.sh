#!/bin/bash

git fetch
git branch -v
git merge origin/main

go run .
