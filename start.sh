#!/bin/sh

set -e

echo "start the app"
exec "$@"   # takes all param and exec
