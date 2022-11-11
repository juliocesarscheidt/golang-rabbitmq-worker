#!/bin/bash

sleep 10

while true; do
  sleep 5
  rabbitmq-diagnostics -q ping > /dev/null 2>&1
  test $? == 0 && break
done

# setup
rabbitmqadmin declare --vhost="${RABBITMQ_DEFAULT_VHOST}" --user "${RABBITMQ_DEFAULT_USER}" --password "${RABBITMQ_DEFAULT_PASS}" queue name="orders" auto_delete=false durable=true
rabbitmqadmin declare --vhost="${RABBITMQ_DEFAULT_VHOST}" --user "${RABBITMQ_DEFAULT_USER}" --password "${RABBITMQ_DEFAULT_PASS}" queue name="orders_error" auto_delete=false durable=true
