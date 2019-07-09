#!/bin/sh
cd /go/bin
pwd
ls -lah
/go/bin/sql-migrate up -env ${ENV}
/go/bin/serverd
