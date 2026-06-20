#!/bin/bash

# Check if a file path is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <path_to_m3u_file>"
    exit 1
fi

INPUT_FILE="$1"

# Check if the input file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: Input file '$INPUT_FILE' not found."
    exit 1
fi

# Create a temporary file for the modified content
TEMP_FILE=$(mktemp)

# Process the input file
while IFS= read -r line; do
    # If the line starts with '#EXTM3U', print it as is
    if [[ "$line" == "#EXTM3U"* ]]; then
        echo "$line" >> "$TEMP_FILE"
        # If the line contains a file path, change the extension to .mp3
        elif [[ "$line" == *.* ]]; then
        # Get the directory and filename without extension
        dirname=$(dirname "$line")
        filename=$(basename "$line" | sed 's/\.[^.]*$//')
        # Construct the new path with .mp3 extension
        echo "$dirname/$filename.mp3" >> "$TEMP_FILE"
        # Otherwise, print the line as is
    else
        echo "$line" >> "$TEMP_FILE"
    fi
done < "$INPUT_FILE"

# Overwrite the original file with the modified content
mv "$TEMP_FILE" "$INPUT_FILE"

echo "M3U playlist '$INPUT_FILE' updated with .m4a and .flac extensions changed to .mp3."
