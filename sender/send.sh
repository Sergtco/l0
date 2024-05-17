#!/bin/bash

FILE=$(cat model.json)

nats pub ORDERS.new "${FILE}"
