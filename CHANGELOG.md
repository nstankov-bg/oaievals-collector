# Changelog

All notable changes to the OAIEvals Collector project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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

[Unreleased]: https://github.com/oaievals-collector/seed/compare/0.0.3-beta1...HEAD
[0.0.3-beta1]: https://github.com/oaievals-collector/seed/compare/0.0.2-beta1...0.0.3-beta1
[0.0.2-beta1]: https://github.com/oaievals-collector/seed/compare/0.0.1-beta1...0.0.2-beta1
[0.0.1-beta1]: https://github.com/oaievals-collector/seed/releases/tag/0.0.1-beta1
