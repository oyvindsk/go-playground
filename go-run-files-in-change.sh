#!/bin/bash

# go run .go in the curent folder, when any .go file in the curent folder is updated
# fswatch can watch sub dirs, exclude, include etc etc.
# See: https://emcrisostomo.github.io/fswatch/ ,  https://aur.archlinux.org/packages/fswatch

fswatch -o --event Updated *.go |  xargs -n2 go run *.go -

