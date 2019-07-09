#!/bin/sh
cd /go/bin
/go/bin/sql-migrate up -env ${ENV}
/go/bin/serverd
