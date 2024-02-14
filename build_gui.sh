#!/usr/bin/env bash

function suffix() {
    if [[ "$OSTYPE" == "msys" ]]; then
        echo "$1.exe"
    else
        echo "$1"
    fi
}

cd ./gui || exit 1
fyne build -o "$(suffix ./middleman-gui)"
mv "$(suffix ./middleman-gui)" ../
