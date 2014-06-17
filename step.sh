#!/bin/bash

bundle install
ruby ./hipchat.rb
exit $?