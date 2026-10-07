#!/bin/sh
set -eu

for name in .ssh .gitconfig; do
    source="/tmp/${name}"
    target="${HOME}/${name}"

    if { [ -e "$source" ] || [ -L "$source" ]; } \
       && [ ! -e "$target" ] \
       && [ ! -L "$target" ]; then
        ln -s "$source" "$target"
    fi
done

exec sleep infinity
