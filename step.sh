#!/bin/bash

curl http://curl.haxx.se/ca/cacert.pem > $HOME/cacert.pem
export SSL_CERT_FILE=$HOME/cacert.pem
bundle install --path vendor/bundle
ruby ./hipchat.rb