#!/bin/sh

go run cmd/cmd.go \
    --width=2048 \
    --height=2048 \
    --out ./img/colored.jpg \
    --mode pixel \
    --scale 1 \
    --threshold 1000 \
    --workers 8 \
    --iter 800 \
    --zoom 20000000 \
    --offsetX 0.7 \
    --offsetY 0.291 \
    --hue 100 \
    --quality 75
