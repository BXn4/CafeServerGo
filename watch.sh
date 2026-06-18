#!/bin/bash

PID=$(pidof cafego)

if [ -n "$PID" ]; then
  echo "Process: cafego (PID: $PID)"
  echo "---------------------------------------------"

  CPU=$(ps -p $PID -o %cpu=)
  MEM_RSS=$(ps -p $PID -o rss=)
  MEM_VSZ=$(ps -p $PID -o vsz=)
  MEM_RSS_MB=$(echo "scale=2; $MEM_RSS/1024" | bc)
  MEM_VSZ_MB=$(echo "scale=2; $MEM_VSZ/1024" | bc)
  THREADS=$(ps -p $PID -o nlwp=)
  CMD=$(ps -p $PID -o cmd=)

  echo "CPU: $CPU%"
  echo "Memory: RSS ${MEM_RSS_MB} MB | VSZ ${MEM_VSZ_MB} MB"
  echo "Threads: $THREADS"
  echo "Command: $CMD"

  echo ""
  echo "Network traffic (system eth0):"
  RX=$(grep 'eth0:' /proc/net/dev | awk -F':' '{gsub(/ /,"",$2); print $2}' | awk '{print $1}')
  TX=$(grep 'eth0:' /proc/net/dev | awk -F':' '{gsub(/ /,"",$2); print $2}' | awk '{print $9}')
  echo "RX: $RX bytes | TX: $TX bytes"
else
  echo "cafego process not running"
fi
