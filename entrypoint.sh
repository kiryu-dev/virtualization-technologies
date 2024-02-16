#!/bin/sh

if [ "$(id -u)" != 0 ]; then
  echo 'container must start as root' >&2
  exit 1
fi

exec gosu kira "$@"