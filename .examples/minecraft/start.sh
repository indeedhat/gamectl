#!/bin/bash

if [[ ! -f ".examples/minecraft/.mc.lock" ]]; then
    echo $(date '+%s') > .examples/minecraft/.mc.lock
else
    exit 1
fi
