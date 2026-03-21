#!/usr/bin/env bash

set -x

minerva='http://localhost:9080'
scan="$(nmap -RF -oX - 192.168.1.0/24)"
curl -qs --data "${scan}" "${minerva}/hosts/"
