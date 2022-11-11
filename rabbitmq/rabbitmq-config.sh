#!/bin/bash

sleep 10

while true; do
  sleep 5
  rabbitmq-diagnostics -q ping > /dev/null 2>&1
  test $? == 0 && break
done
