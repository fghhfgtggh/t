#!/bin/bash
./or &
sleep 2
exec socat TCP-LISTEN:${PORT},fork TCP:127.0.0.1:6324
