#!/bin/bash

# Build the project
echo "Building the project..."

# Check if the build directory exists
if [ ! -d "build" ]; then
    echo "Creating the build directory..."
    mkdir build
fi

# Remove the existing build

# Check if the build directory is empty
if [ "$(ls -A build)" ]; then
    echo "Removing the existing build..."
    rm -rf build/*
fi

echo "Copying raylib library..."

# Copy the raylib library
cp -r libs/raylib/windows/raylib.dll build/

echo "Building the project..."

# Build the project
go build -o build/fenrir cmd/fenrir/main.go

# Check if the build was successful
if [ $? -eq 0 ]; then
    echo "Build successful!"
    exit 0
else
    echo "Build failed!"
    exit 1
fi
