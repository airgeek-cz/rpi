#!/bin/sh
set -eu

SERIAL=/dev/ttyS0
NRST_GPIO=26

# Reset Airgeek
raspi-gpio set "$NRST_GPIO" op dl
sleep 1
raspi-gpio set "$NRST_GPIO" op dh

# Configure serial port and open as stdin
stty -F "$SERIAL" 115200 -igncr
exec <"$SERIAL"

# Launch the sidecar
exec go run .
