#!/bin/sh

echo 'Running'
export GO_ENV=production
exec "$@"
