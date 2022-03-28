# sample-netflow
Sample Project for exposing net flow stats

Start Postgress DB
docker-compose up -d --build --force-recreate

Execute the following to create NetFlow Table for the postgressDB

docker exec -it netflowdb /bin/sh

psql -U postgres

CREATE TABLE netflow (
  src_app    VARCHAR(250) NOT NULL,
  dest_app   VARCHAR(250) NOT NULL,
  vpc_id     VARCHAR(250) NOT NULL,
  bytes_tx   NUMERIC,
  bytes_rx   NUMERIC,
  hour       NUMERIC
);

#Launch the app
go run main.go
