#!/bin/bash

if [[ -f ".examples/minecraft/.mc.lock" ]]; then
    rm .examples/minecraft/.mc.lock
else
    exit 1
fi
