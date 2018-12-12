#!/bin/sh

# Check to see if we need to run the seed process
if [ -n "$SEED" ]
then
   echo 'Initializing and seeding the database'
   # Retry migration in case the db is still starting up.
   n=0
   while [ $n -lt 10 ]
   do
       cms_authz_linux migrate && cms_authz_linux task db:seed && break
       n=`expr $n + 1`
       echo 'Retrying the migration command.'
       sleep 10
   done
fi

echo 'Starting server in prod mode'
exec "$@"
