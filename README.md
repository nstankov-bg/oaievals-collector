# OAIEvals Collector

[![Keep a Changelog](https://img.shields.io/badge/changelog-Keep%20a%20Changelog-%23E05735)](CHANGELOG.md)
[![GitHub Release](https://img.shields.io/github/v/release/nstankov-bg/oaievals-collector)](https://github.com/nstankov-bg/oaievals-collector/releases)
[![Go Reference](https://pkg.go.dev/badge/nstankov-bg/oaievals-collector.svg)](https://pkg.go.dev/github.com/nstankov-bg/oaievals-collector)
[![go.mod](https://img.shields.io/github/go-mod/go-version/nstankov-bg/oaievals-collector)](go.mod)
[![LICENSE](https://img.shields.io/github/license/nstankov-bg/oaievals-collector)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/nstankov-bg/oaievals-collector)](https://goreportcard.com/report/github.com/nstankov-bg/oaievals-collector)

[![Docker Pulls](https://badgen.net/docker/pulls/nikoogle/oaievals-collector?icon=docker&label=pulls)](https://hub.docker.com/repository/docker/nikoogle/oaievals-collector)
[![Docker Stars](https://badgen.net/docker/stars/nikoogle/oaievals-collector?icon=docker&label=stars)](https://hub.docker.com/repository/docker/nikoogle/oaievals-collector)
[![Docker Image Size](https://badgen.net/docker/size/nikoogle/oaievals-collector?icon=docker&label=image%20size)](https://hub.docker.com/repository/docker/nikoogle/oaievals-collector)


⭐ `Star` this repository if you find it valuable and worth maintaining.

👁 `Watch` this repository to get notified about new releases, issues, etc.

## Stats

![Alt](https://repobeats.axiom.co/api/embed/8a8f307f96fa2ebb08879a294ad1ddc909ad9619.svg "Repobeats analytics image")

## Introduction

The OAIEvals Collector is a Go application specifically designed to collect and store raw evaluation metrics. It provides an HTTP handler & FileSystem Watcher for receiving metric data, a set of utilities for event handling and monitoring, and supports containerized deployment with Docker for ease of use and scalability.

The application integrates seamlessly with InfluxDB for robust and efficient storage of numeric-based time series data, Loki for storing and retrieving log data, providing context and qualitative information around numeric metrics, and now with TimescaleDB for storing and managing event data. All three systems can be used independently or together, depending on the nature of the metrics and the requirements of your system.

The latest addition to the list of integrations is Elasticsearch 8. Elasticsearch is a real-time distributed search and analytics engine built on Apache Lucene. It provides scalable search capabilities across vast amounts of data and complements our existing storage solutions. It allows for complex queries and aggregations, making it suitable for diverse use-cases, such as log and event data analysis. With this addition, alongside Kafka and MongoDB, we now provide a comprehensive data ingestion and storage solution. Kafka enables efficient ingestion and processing of high volumes of real-time events, suitable for systems handling thousands of events per second. MongoDB offers flexible, document-based data modeling, capable of handling a wide variety of data types.

## Demo:

**Elasticsearch x Kibana:**

| Elasticsearch Test |
| :---: |
| *Basic Visualization via Kibana*
| ![Screenshot 2023-07-22 at 10 33 28 AM](https://github.com/openai/evals/assets/27363885/e03a4cf7-49f4-4207-a1b7-58a5ccef9b1c)

**InfluxDB:**

| InfluxDB Test 1 | InfluxDB Test 2 |
| :---: | :---: |
| ![InfluxDB Image 1](https://github.com/openai/evals/assets/27363885/f2359bce-5af2-49c6-a4dd-66e362ece63d) | ![InfluxDB Image 2](https://github.com/openai/evals/assets/27363885/be3c7361-2601-417d-a311-96c09da954c9) |

**Grafana:**

| Grafana Test |
| :---: |
| *Basic Visualization via InfluxDB & TimeScaleDB*
| ![Grafana Image](https://github.com/nstankov-bg/oaievals-collector/assets/27363885/cd119b1a-939c-4f2d-b141-d26e83784cbc) 

**Kafka:**

| Kafka Test |
| :---: |
| *Basic Visualization via Kafka UI*
| ![Kafka Image](https://github.com/nstankov-bg/oaievals-collector/assets/27363885/e8075f06-b628-4773-99d9-a032e28f2472) 

## Prerequisites

Before you can run the OAIEvals Collector, you'll need to have the following installed:

- [Go](https://golang.org/dl/) (version 1.20 or later)
- [Docker](https://www.docker.com/products/docker-desktop) (version 20.10 or later)
- [Docker Compose](https://docs.docker.com/compose/install/) (version 1.28 or later)

Please ensure you have these prerequisites installed and properly configured before proceeding with the setup of the OAIEvals Collector.

## Usage

The OAIEvals Collector is designed to be deployed as a containerized application, with its services orchestrated by Docker Compose. The application can be configured using environment variables.

**Preparation:**

1. Clone the repository and navigate to the project directory.

**Start Services:**

1. Run `docker-compose up` to start the OAIEvals Collector, Kafka,Elasticsearch, InfluxDB, Loki, TimescaleDB, MongoDB. Docker will pull the images (if not already present) and build the OAIEvals Collector image.

**InfluxDB, Loki, TimescaleDB, Kafka, and MongoDB Setup:**

1. While the services are spinning up, navigate to your InfluxDB instance and generate an authentication token (for InfluxDB) and/or get the connection strings. These tokens and strings will be used by the OAIEvals Collector to connect and interact with the databases.

**Configuration Setup:**

1. Once you've obtained the tokens and connection strings, stop the running Docker Compose services (using CTRL+C or `docker-compose down` command).
2. Open the `.env` file (create one based on the provided `.env.example` if it does not exist), and replace `your_token_here` in `INFLUXDB_TOKEN=your_token_here`, `LOKI_URL=your_loki_url`,`KAFKA_BOOTSTRAP_SERVERS=your_kafka_bootstrap_servers`, `TIMESCALEDB_HOST=your_timescaledb_host`, `MONGODB_URI=your_mongodb_uri`, `ES_ADDRESS=your_elastic_address` and others with the tokens or endpoints obtained from InfluxDB, Loki, TimescaleDB, Kafka, and MongoDB respectively.

## Dependencies

The OAIEvals Collector is compatible with various storage systems for metrics collection, and it allows the user to select one or more collection points according to their needs. The supported storage systems include [InfluxDB](https://www.influxdata.com/), [Loki](https://grafana.com/oss/loki/), [TimescaleDB](https://www.timescale.com/), [Kafka](https://kafka.apache.org/), [MongoDB](https://www.mongodb.com/), and now, [Elasticsearch 8](https://www.elastic.co/). Ensure you have instances of these systems available for the application to connect to if you choose to use them.

## Build, Run, Test

### Running with Docker Compose

You can easily run the application along with the desired data storage service (InfluxDB, Loki, TimescaleDB, Kafka, or MongoDB) using Docker Compose. For each data storage service, a dedicated Docker Compose file and corresponding Make commands are provided.

To start the application with a specific service, view the logs, or stop the service, use the appropriate Make command. Here's how you do it for each service:

- **TimescaleDB**: `make up-timescale`, `make logs-timescale`, `make down-timescale`
- **MongoDB**: `make up-mongo`, `make logs-mongo`, `make down-mongo`
- **InfluxDB**: `make up-influx`, `make logs-influx`, `make down-influx`
- **Loki**: `make up-loki`, `make logs-loki`, `make down-loki`
- **Kafka**: `make up-kafka`, `make logs-kafka`, `make down-kafka`
- **Elastic**: `make up-elastic`, `make logs-elastic`, `make down-elastic`

For example, `make up-timescale` starts the application with TimescaleDB, `make logs-timescale` displays its logs, and `make down-timescale` stops the service.

### Locally

1. Build the application: `go build -o oaievals-collector .`
1. Run the application: `./oaievals-collector`

**Test the Application:**

You can test the OAIEvals Collector by sending a POST request with a payload representing an event. For instance:

```shell
curl -X POST \
     -H "Content-Type: application/json" \
     -d '[
         {
           "run_id": "2307080128125Q6U7IFP",
           "event_id": 3,
           "sample_id": "abstract-causal-reasoning-text.dev.484",
           "type": "match",
           "data": {
             "correct": false,
             "expected": "ANSWER: off",
             "picked": "undetermined.",
             "sampled": "undetermined."
           },
           "created_by": "",
           "created_at": "2023-07-08 01:28:13.704853+00:00"
         }
       ]' http://localhost:8081/events
```


## Contributing

Feel free to create an issue or propose a pull request. Follow the [Code of Conduct](CODE_OF_CONDUCT.md).
