# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-06-09

### Added

* support for using the latest remote `release.json` schema if no `--schema` file provided
* releases for mac and windows
* this changelog!

### Changed

* `--schema` is now optional
* bumped `santhosh-tekuri/jsonschema` from `v5` to `v6`
* bumped Go from `1.21` to `1.24`

## [0.2.0] - 2024-05-05

### Added

* `.jsonl` support. The `--in` parameter now accepts `.jsonl` input.
    - each line in the `.jsonl` file will be validated as an individual `release.json` file.

### Changed

* the `--in` parameter must now have either a `.json` or `.jsonl` extension.
* the `--schema` parameter must now have a `.json` extension.

## [0.1.0] - 2024-05-05

### Added

* initial release with barest minimal functionality
* `--in` to specify input file
* `--schema` to specify schema to validate against
* `--help` to print the options
* `--version` to print version of the tool
* `./manage.sh` with 'build', 'release' and 'clean' type commands

## [Unreleased]

### Added

### Changed

### Deprecated

### Removed

### Fixed

### Security
