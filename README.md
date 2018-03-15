# wrkr-t3chfest: Columnar databases and Clickhouse T3CHFEST 2018.

This is the a description of the how to run the demo for my talk about columnar databases and Clickhouse on T3CHFEST 2018. 

The slides can be found here https://es.slideshare.net/FcoJavierSanzOlivera/click-houset3chfest (in Spanish)
and the video of talk 
https://www.youtube.com/watch?v=_oZKi5v951Q (also in Spanish)

## Docker Compose

This example includes a docker compose that allows you to run the next services:

- A Clickhouse database instance
- A MariaDB database instance
- All the services needed to run a Redash (https://redash.io/) (Redis, Posgres, Ngix, etc)

You need to install docker on your machine in order to run this demo. Once you have that installed you need to run the next commands:

 `docker-compose up -d`

This will boot up all the needed services. Running a:

 `docker ps` 

you will see something like:

```
CONTAINER ID        IMAGE                                COMMAND                  CREATED                  STATUS              PORTS                                                      NAMES
afd034bd71d6        golang:1.9.4-alpine3.7               "sh"                     Less than a second ago   Up 4 seconds                                                                   wrkrt3chfest_gobuild_1
bea975585a4f        redash/nginx:latest                  "nginx -g 'daemon of…"   6 minutes ago            Up 7 minutes        0.0.0.0:80->80/tcp, 443/tcp                                wrkrt3chfest_nginx_1
88a4ec46a8d5        redash/redash:4.0.0-beta.b3690       "/app/bin/docker-ent…"   6 minutes ago            Up 7 minutes        0.0.0.0:5000->5000/tcp                                     wrkrt3chfest_server_1
ad2df38414bd        mariadb:10.3                         "docker-entrypoint.s…"   6 minutes ago            Up 7 minutes        0.0.0.0:3306->3306/tcp                                     wrkrt3chfest_mariadb_1
d9b9db13b53e        postgres:9.5.6-alpine                "docker-entrypoint.s…"   6 minutes ago            Up 7 minutes        0.0.0.0:5432->5432/tcp                                     wrkrt3chfest_postgres_1
d2c78688acc2        redis:3.0-alpine                     "docker-entrypoint.s…"   6 minutes ago            Up 7 minutes        6379/tcp                                                   wrkrt3chfest_redis_1
6b90950d21d6        redash/redash:4.0.0-beta.b3690       "/app/bin/docker-ent…"   6 minutes ago            Up 7 minutes        5000/tcp                                                   wrkrt3chfest_worker_1
40d92e7a0591        yandex/clickhouse-server:1.1.54327   "/bin/sh -c 'exec /u…"   6 minutes ago            Up 7 minutes        0.0.0.0:8123->8123/tcp, 0.0.0.0:9000->9000/tcp, 9009/tcp   wrkrt3chfest_clickhouse_1
```

Now you need to create the tables for Redash. So execute the command

`docker exec -it wrkrt3chfest_server_1 bash`

And inside the shell on the container run

`edash@3499c2b2e3c1:/app$ ./manage.py database create_tables`

`edash@3499c2b2e3c1:/app$ exit`

No you can open localhost on your browser and you will need to register on Redash. Once done, the next step is to add a Data source. Click your profile(up right)->data sources
and then `+`. Then the parameters you need to add are:

The default password is '' but Redash does not handle that very well so you can put `default`. Click `Test Connection` and if the green Widget appears you are ready to use your
Clickhouse. 

##Populating Data

On my demo on the talk I used Backblaze Hard Drive Data and Stats https://www.backblaze.com/b2/hard-drive-test-data.html. 
It is a dataset of information (mainly failures rates) of their
own real systems. They provide several quarters of info on CSV files. 
You could use different methods included both on Clickhouse and on MariaDB to populate CSV data. 
But in this case I provide a Golang Script that is ready to load all data from 2017-01-01 to 2017-09-30. I case you wanted to load
different dates there is a `for loop` at the end of the script that can be modified:

```
for curDate := date.New(2017, 1, 1); date.New(2017, 9, 30).After(curDate); curDate = curDate.AddDays(1) {
```

The CVS files are expected to be on the folder `data_2017/`. In case you want a different folder change that on the `.yaml` files include on the
project. There is two because this script allows to populate both MariaDB and Clickhouse. On MariaDB the database `t3chfest` must be created 
beforehand. But again that can be change that on the `.yaml` file. 

You don't need to install Golang on your system to run the script. The docker-compose includes a golang build image as well. Run:

`docker exec -it wrkrt3chfest_gobuild_1 sh`

`/go # cd src/wrkr-t3chfest/`

And then either run 

`go run populate.go clickhouse``

for populating data on Clickhouse or 

`go run populate.go mariadb`

for MariaDB

No you can run queries on both DBs and compare results. On file the `queries.sql` you can see some examples.

Enjoy!!!
Javier



