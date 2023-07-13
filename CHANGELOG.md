# Changelog

All notable changes to the OAIEvals Collector project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- MongoDB support: The OAIEvals Collector now includes MongoDB as a target for collecting and storing event data. With the new addition, the collector can connect, interact, and write event data to MongoDB databases. The MongoDB support allows for flexible, document-based data modeling and can handle a wide variety of data types. Required environment variables: `MONGODB_URI`, `MONGODB_DATABASE`, `MONGODB_COLLECTION`.

## [0.0.4-beta1] - 2023-07-13

### Added
- Kafka support: We've added Kafka as a target for collecting and storing event data. The OAIEvals Collector can now connect, interact with Kafka, and write messages to Kafka topics at a high rate, thus handling large volumes of event data. The Kafka support includes the buffering of messages for faster ingestion and asynchronous writing. Required environment variable: `KAFKA_BOOTSTRAP_SERVERS`.

## [0.0.3-beta1] - 2023-07-11

### Added
- TimescaleDB support: We've added TimescaleDB as a target for collecting and storing event data. The OAIEvals Collector can now connect and interact with TimescaleDB, providing flexibility in handling different types of event data and supporting time-series analyses.

## [0.0.2-beta1] - 2023-07-11

### Added
- Loki support: We've added Loki as a target for collecting and storing log data. The OAIEvals Collector can now connect and interact with Loki, complementing the existing InfluxDB integration. This offers more flexibility in handling different types of metrics, providing context and qualitative information around numeric metrics.

## [0.0.1-beta1] - 2023-06-08

### Added
- Initial release: Created OAIEvals Collector, a Go application designed to collect and store raw evaluation metrics.
- HTTP Handler & FileSystem Watcher: These components are responsible for receiving metric data.
- Integration with InfluxDB: For robust and efficient data storage, we've integrated the collector with InfluxDB to handle numeric-based time series data.
- Containerized Deployment: The application supports deployment with Docker and orchestration with Docker Compose, making it easy to use and scalable.

[Unreleased]: https://github.com/oaievals-collector/seed/compare/0.0.4-beta1...HEAD
[0.0.4-beta1]: https://github.com/oaievals-collector/seed/compare/0.0.3-beta1...0.0.4-beta1
[0.0.3-beta1]: https://github.com/oaievals-collector/seed/compare/0.0.2-beta1...0.0.3-beta1
[0.0.2-beta1]: https://github.com/oaievals-collector/seed/compare/0.0.1-beta1...0.0.2-beta1
[0.0.1-beta1]: https://github.com/oaievals-collector/seed/releases/tag/0.0.1-beta1