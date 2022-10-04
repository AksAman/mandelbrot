#!/bin/sh

# Install
# go install github.com/cosmtrek/air@latest

# go run main.go \
# ./main \
air -- \
    --mode pixel \
    --scale 2 \
    --threshold 128 \
    --workers 8 \
    --iter 1000
    --zoom 2 \
    --offsetX 0.7435 \
    --offsetY -0.1315 \
    --hue 200 \
