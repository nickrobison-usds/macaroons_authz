#!/bin/sh

export GO_ENV=production

# Check to see if we need to run the seed process
if [ -n "$SEED" ]
then
   echo 'Initializing and seeding the database'
   # Run the migration
   cms_authz_linux migrate

   # Seed the data
   cms_authz_linux task db:seed
fi

echo 'Starting server in prod mode'
exec "$@"
