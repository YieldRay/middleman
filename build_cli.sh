#!/usr/bin/env bash

function suffix() {
    if [[ "$OSTYPE" == "msys" ]]; then
        echo "$1.exe"
    else
        echo "$1"
    fi
}

go build -o "$(suffix middleman)" ./cli


