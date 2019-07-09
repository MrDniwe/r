#!/bin/sh
cd /go/bin
sql-migrate up -env ${ENV}
/go/bin/serverd
