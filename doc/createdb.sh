#!/usr/bin/env bash

sqlite3 "$1.db" < "${2}"
