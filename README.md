# sample-netflow
Sample Project for exposing net flow stats

Start Postgress DB
docker-compose up -d --build --force-recreate

Execute the following to create NetFlow Table for the postgressDB

docker exec -it netflowdb /bin/sh

psql -U postgres

'[{"src_app": "foo", "dest_app": "bar", "vpc_id": "vpc-0", "bytes_tx":
100, "bytes_rx": 500, "hour": 1}]'


CREATE TABLE netflow (
  id SERIAL,
  src_app    VARCHAR(250) NOT NULL,
  dest_app   VARCHAR(250) NOT NULL,
  vpc_id     VARCHAR(250) NOT NULL,
  bytes_tx   NUMERIC,
  bytes_rx   NUMERIC,
  hour       NUMERIC,
  PRIMARY KEY (src_app, dest_app, vpc_id)
);


