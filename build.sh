#!/bin/bash

git add -A

git commit -m "Update"

git push

version=$1

if [ -z "$1" ]; then
  echo "Version is not provided. Exiting..."
  exit 1
fi


echo "Version: $version"

# Build the project
echo "Building the project..."

git tag $version -m "Release $version"
git push origin $version