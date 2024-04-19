# Bisquitt-PSK

### Simple PSK service / library as API for loading and serving PSK (pre-shared keys) to [BISQUITT](https://github.com/energostack/bisquitt) service

**Can be used as a standalone service or just as example of implementation**

## Usage
Change environment variables in `docker-compose.yml` to your own values.

Run the following command to start the service:
```bash
docker-compose up
```

Currently service supports CSV, SQLITE and Kafka as data sources.
To change the data source, change the `DATA_SOURCE` environment variable in `docker-compose.yml`
to one of the following values: `csv`, `sqlite`, `kafka`.

Don't forget to change the corresponding
`FILE_PATH` for `csv` and `sqlite` data sources or KAFKA environment variables for Kafka data source.

**Change `BASIC_AUTH_USERNAME` and `BASIC_AUTH_PASSWORD` to ensure security.**

## API

##### Service provides OpenAPI documentation at `/docs/index.html` endpoint.

Header `Authorization: Basic base64(username:password)` is required for all requests (except docs).

### GET /clients/{client_id}
Returns a PSK for a specific client as a JSON object.
```bash
{"client":"bisquitt","psk":"cHNr"}
```
