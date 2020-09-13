#!/bin/sh
set -x
echo "Hot reload script starts :)"
reflex -r "(\.go$|go\.mod$)" -s -- sh -c "go run main.go"