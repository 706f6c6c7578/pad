#!/bin/bash

# Input the file size in bytes
read -p "Enter the file size in bytes: " file_size

# Input the block size in bytes
read -p "Enter the block size in bytes: " block_size

# Calculate the number of blocks
num_blocks=$(( (file_size + block_size - 1) / block_size ))

# Calculate the new file size
new_file_size=$(( num_blocks * block_size ))

echo "New file size: $new_file_size bytes"

