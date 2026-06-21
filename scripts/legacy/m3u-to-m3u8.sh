#!/bin/bash

# Default values for flags
FORCE_OVERWRITE=false
REMOVE_ORIGINAL_FORCE=false # Flag to force removal of original without prompt

# Input file variable
INPUT_FILE=""

# Parse command line arguments
while [[ "$#" -gt 0 ]]; do
    case "$1" in
        -f|--force)
            FORCE_OVERWRITE=true
        ;;
        -r|--remove-original) # This flag now forces removal
            REMOVE_ORIGINAL_FORCE=true
        ;;
        -*) # Any other flag starting with -
            echo "Error: Unknown flag '$1'"
            echo "Usage: $0 [-f|--force] [-r|--remove-original] <input.m3u/.m3u8>"
            exit 1
        ;;
        *) # Positional argument (assumed to be the input file)
            if [ -z "$INPUT_FILE" ]; then
                INPUT_FILE="$1"
            else
                echo "Error: Too many input files specified. Only one expected."
                echo "Usage: $0 [-f|--force] [-r|--remove-original] <input.m3u/.m3u8>"
                exit 1
            fi
        ;;
    esac
    shift # Move to the next argument
done

if [ -z "$INPUT_FILE" ]; then
    echo "Usage: $0 [-f|--force] [-r|--remove-original] <input.m3u/.m3u8>"
    exit 1
fi

OUTPUT_FILE=""

# Determine output filename and handle .m3u8 input
if [[ "$INPUT_FILE" =~ \.m3u$ ]]; then
    OUTPUT_FILE="${INPUT_FILE%.m3u}.m3u8"
    elif [[ "$INPUT_FILE" =~ \.m3u8$ ]]; then
    OUTPUT_FILE="$INPUT_FILE"
else
    echo "Error: Input file must be .m3u or .m3u8"
    exit 1
fi

if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: Input file '$INPUT_FILE' not found."
    exit 1
fi

if [ "$FORCE_OVERWRITE" = false ]; then # Only prompt if not forced
    if [ -f "$OUTPUT_FILE" ] && [ "$INPUT_FILE" != "$OUTPUT_FILE" ]; then
        read -p "Warning: Output file '$OUTPUT_FILE' already exists. Overwrite? (y/N): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            echo "Operation cancelled."
            exit 0
        fi
        elif [ -f "$OUTPUT_FILE" ] && [ "$INPUT_FILE" = "$OUTPUT_FILE" ]; then
        # If input and output are the same (m3u8 -> m3u8), check if it needs fixing
        if grep -q "#EXTM3U" "$INPUT_FILE" && ! grep -q "#EXTINF" "$INPUT_FILE"; then
            read -p "Warning: '$INPUT_FILE' is an M3U8 file missing #EXTINF tags. Overwrite to fix? (y/N): " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                echo "Operation cancelled."
                exit 0
            fi
        else
            echo "Info: '$INPUT_FILE' is an M3U8 file and appears to have #EXTINF tags. No action needed."
            exit 0
        fi
    fi
fi

# Check for ffprobe
if ! command -v ffprobe &> /dev/null; then
    echo "Error: ffprobe is not installed. Please install FFmpeg to proceed."
    exit 1
fi

# Determine the directory of the M3U/M3U8 file for relative paths
FILE_DIR="$(dirname "$INPUT_FILE")"

# Start M3U8 content with just the basic header
echo "#EXTM3U" > "$OUTPUT_FILE" # This now unconditionally overwrites if past prompts

while IFS= read -r line; do
    if [[ "$line" =~ ^#EXTM3U$ || "$line" =~ ^#EXT-X-VERSION:[0-9]+$ || "$line" =~ ^#EXTINF:.* || "$line" =~ ^#EXT-X-TARGETDURATION:.* || "$line" =~ ^#EXT-X-MEDIA-SEQUENCE:.* || "$line" =~ ^#EXT-X-ENDLIST$ ]]; then
        # Skip all existing M3U/M3U8 headers and EXTINF tags
        continue
        elif [[ "$line" =~ ^[^#].* ]]; then # Media file path
        # Resolve the absolute path of the media file
        MEDIA_PATH="$(realpath "$FILE_DIR/$line")"
        
        if [ -f "$MEDIA_PATH" ]; then
            # Get duration using ffprobe, truncate to an integer
            DURATION=$(ffprobe -v error -show_entries format=duration -of default=noprint_wrappers=1:nokey=1 "$MEDIA_PATH" | cut -d'.' -f1)
            
            if [ -z "$DURATION" ]; then
                echo "Warning: Could not get duration for '$line'. Skipping."
                continue
            fi
            
            # Get Artist and Title from audio file tags
            ARTIST=$(ffprobe -v error -show_entries format_tags=artist -of default=noprint_wrappers=1:nokey=1 "$MEDIA_PATH")
            TITLE=$(ffprobe -v error -show_entries format_tags=title -of default=noprint_wrappers=1:nokey=1 "$MEDIA_PATH")
            
            ARTIST_TITLE=""
            if [ -n "$ARTIST" ] && [ -n "$TITLE" ]; then
                ARTIST_TITLE="${ARTIST} - ${TITLE}"
                elif [ -n "$TITLE" ]; then # Fallback to just title if artist is missing
                ARTIST_TITLE="$TITLE"
            else
                # Fallback to filename parsing if no tags found
                FILENAME=$(basename "$line")
                FILENAME_NO_EXT="${FILENAME%.*}" # Remove extension
                
                if [[ "$FILENAME_NO_EXT" =~ ^(.*)\ -\ (.*)$ ]]; then
                    ARTIST_TITLE="${BASH_REMATCH[1]} - ${BASH_REMATCH[2]}"
                else
                    ARTIST_TITLE="$FILENAME_NO_EXT" # Fallback to raw filename
                fi
            fi
            
            echo "#EXTINF:${DURATION},${ARTIST_TITLE}" >> "$OUTPUT_FILE"
            echo "$line" >> "$OUTPUT_FILE" # Write the original relative path
        else
            echo "Warning: Media file not found at '$MEDIA_PATH' (original path: '$line'). Skipping."
        fi
    fi
done < "$INPUT_FILE"

echo "Processed '$INPUT_FILE' to '$OUTPUT_FILE'"

# Only consider removing if the output file is different from the input file
if [ "$INPUT_FILE" != "$OUTPUT_FILE" ]; then
    if [ "$REMOVE_ORIGINAL_FORCE" = true ]; then
        # If -r flag is present, force remove without asking
        rm "$INPUT_FILE"
        echo "Removed original file '$INPUT_FILE' (forced)."
    else
        # If -r flag is NOT present, prompt the user
        read -p "Remove original file '$INPUT_FILE'? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm "$INPUT_FILE"
            echo "Removed '$INPUT_FILE'."
        else
            echo "Original file retained."
        fi
    fi
else
    echo "Info: Original file is the same as output file, not removing."
fi
