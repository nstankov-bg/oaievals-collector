# OAIEvals Collector

[![Keep a Changelog](https://img.shields.io/badge/changelog-Keep%20a%20Changelog-%23E05735)](CHANGELOG.md)
[![GitHub Release](https://img.shields.io/github/v/release/oaievals-collector/seed)](https://github.com/oaievals-collector/seed/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/oaievals-collector/seed.svg)](https://pkg.go.dev/github.com/oaievals-collector/seed)
[![go.mod](https://img.shields.io/github/go-mod/go-version/oaievals-collector/seed)](go.mod)
[![LICENSE](https://img.shields.io/github/license/oaievals-collector/seed)](LICENSE)
[![Build Status](https://img.shields.io/github/actions/workflow/status/oaievals-collector/seed/build.yml?branch=main)](https://github.com/oaievals-collector/seed/actions?query=workflow%3Abuild+branch%3Amain)
[![Go Report Card](https://goreportcard.com/badge/github.com/oaievals-collector/seed)](https://goreportcard.com/report/github.com/oaievals-collector/seed)
[![Codecov](https://codecov.io/gh/oaievals-collector/seed/branch/main/graph/badge.svg)](https://codecov.io/gh/oaievals-collector/seed)

⭐ `Star` this repository if you find it valuable and worth maintaining.

👁 `Watch` this repository to get notified about new releases, issues, etc.

## Introduction

The OAIEvals Collector is a Go application specifically designed to collect and store raw evaluation metrics. It provides an HTTP handler & FileSystem Watcher for receiving metric data, a set of utilities for event handling and monitoring, integrates seamlessly with InfluxDB for robust and efficient data storage, and supports containerized deployment with Docker for ease of use and scalability. Its purpose is to facilitate the monitoring and evaluation process by providing a one-stop solution for metric collection.

## Prerequisites

Before you can run the OAIEvals Collector, you'll need to have the following installed:

- [Go](https://golang.org/dl/) (version 1.14 or later)
- [Docker](https://www.docker.com/products/docker-desktop) (version 20.10 or later)
- [Docker Compose](https://docs.docker.com/compose/install/) (version 1.28 or later)

Please ensure you have these prerequisites installed and properly configured before proceeding with the setup of the OAIEvals Collector.

## Usage

The OAIEvals Collector is designed to be deployed as a containerized application, with its services orchestrated by Docker Compose. The application can be configured using environment variables.

**Preparation:**

1. Clone the repository and navigate to the project directory.

**Start Services:**

1. Run `docker-compose up` to start the OAIEvals Collector and InfluxDB. Docker will pull the InfluxDB image (if not already present) and build the OAIEvals Collector image.

**InfluxDB Setup:**

1. While the services are spinning up, navigate to your InfluxDB instance and generate an authentication token. This token will be used by the OAIEvals Collector to connect and interact with the InfluxDB.

**Configuration Setup:**

1. Once you've obtained the token, stop the running Docker Compose services (using CTRL+C or `docker-compose down` command).
2. Open the `.env` file (create one based on the provided `.env.example` if it does not exist), and replace `your_token_here` in `INFLUXDB_TOKEN=your_token_here` with the token obtained from InfluxDB.

**Restart Services:**

1. Re-run `docker-compose up`. Now, the OAIEvals Collector should be able to connect to the InfluxDB with the provided token.

**Test the Application:**

You can test the OAIEvals Collector by sending a POST request with a payload representing an event. For instance:

```shell
curl -X POST -H "Content-Type: application/json" -d '
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
}' http://localhost:8080/events
```

## Dependencies

The OAIEvals Collector relies on [InfluxDB](https://www.influxdata.com/) for storing collected metrics. Ensure you have an instance of InfluxDB available for the application to connect to.

## Build and Run

### With Docker Compose

A Docker Compose file is provided for convenience, which will start up both the application and an instance of InfluxDB:

1. Build and start the services: `docker-compose up --build`

### Locally

1. Build the application: `go build -o oaievals-collector .`
1. Run the application: `./oaievals-collector`

## Contributing

Feel free to create an issue or propose a pull request.

Follow the [Code of Conduct](CODE_OF_CONDUCT.md).

```

```
# oaievals-collector
