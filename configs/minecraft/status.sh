#!/bin/bash

STATUS="false"
PLAYERS=0
UPTIME=0

if [[ -f "configs/minecraft/.mc.lock" ]]; then
    STATUS="true"
    PLAYERS=1
    
    startTime=$(cat .examples/minecraft/.mc.lock)
    now=$(date '+%s')
    let UPTIME=${now}-${startTime}
fi

echo "{\"online\":${STATUS},\"players\":${PLAYERS},\"max_players\":10,\"uptime\":${UPTIME}}"
