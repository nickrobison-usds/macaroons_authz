#!/bin/sh

exec /usr/bin/cfssl "serve" "-ca" "/cfssl/keys/ca.pem" "-ca-key" "/cfssl/keys/ca-key.pem" "-config" "/cfssl/config/config_ca.json" "-address" "0.0.0.0"
