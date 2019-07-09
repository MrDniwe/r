#!/bin/sh
cd /go/bin
ls -lah
/go/bin/sql-migrate up -env ${ENV}
/go/bin/serverd
