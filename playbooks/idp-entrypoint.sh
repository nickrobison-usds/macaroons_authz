#!/bin/sh

if [ -n "$SEED" ]
then
    echo 'Setting up databases'
    yarn install
    rake db:create
    rake db:environment:set
    rake db:reset
    rake db:environment:set
    rake dev:prime
    rake db:create RAILS_ENV=test
    rake db:reset RAILS_ENV=test
fi

exec "$@"
