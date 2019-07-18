# The R. website backend
## Installation
Release version is oriented to work in tiny *alpina* docker container. Just follow .Dockerfile to build one.

Sample backend + db combination you can see in *docker-compose.yml* in this repository. 

Environment config values are described below:

* `PG_HOST` - Postgres host address, default is `localhost`
* `PG_PORT` - Postgres host port, default is `5434`
* `PG_USER` - Postgres user, default is `development`
* `PG_PASSWORD` - Postgres user's password, default is `development`
* `PG_DATABASE` - Postgres database name, default is `development`
* `S3_URI_PREFIX` - AWS bucket URI where do we store our content images, default is `https://r57.s3.eu-central-1.amazonaws.com`
* `SITE_PAGE_AMOUNT` - Number of articles on a page of article list
