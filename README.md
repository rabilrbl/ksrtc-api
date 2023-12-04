# ksrtc-api

**ksrtc-api** is a simple Go application that provides an API for accessing information related to Karnataka State Road Transport Corporation (KSRTC) buses. This API allows you to retrieve information about bus routes, schedules, and more. It utilizes a caching mechanism to improve performance and responsiveness.

## Table of Contents
- [ksrtc-api](#ksrtc-api)
  - [Table of Contents](#table-of-contents)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
  - [Usage](#usage)
    - [Endpoints](#endpoints)
    - [Examples](#examples)
      - [Home Endpoint](#home-endpoint)
      - [Retrieve All Bus Routes](#retrieve-all-bus-routes)
      - [Retrieve Specific Bus Routes](#retrieve-specific-bus-routes)
  - [Caching](#caching)
  - [Contributing](#contributing)
  - [License](#license)

## Getting Started

### Download the binary

1. Download latest released binary from [releases](https://github.com/rabilrbl/ksrtc-api/releases/latest).
2. On unix based systems, provide executable permissions.
   ```bash
   chmod +x ksrtc-api
   ```
3. Execute the binary and access the API at `http://localhost:8080/`.

## Development

### Prerequisites

Before you can run the `ksrtc-api` project, make sure you have the following installed on your system:

- Go: You can download and install Go from the official [Go website](https://golang.org/doc/install).

### Installation

1. Clone the `ksrtc-api` repository to your local machine:

   ```bash
   git clone https://github.com/rabilrbl/ksrtc-api.git
   cd ksrtc-api
   ```

2. Build and run the application:

   ```bash
   go run main.go
   ```

The API should now be running locally on port 8080, or you can specify a different port by setting the `PORT` environment variable.

## Usage

### Endpoints

The `ksrtc-api` provides the following endpoints:

- `/` - Home endpoint to check if the API is up and running.
- `/all` - Retrieve information about all journey place data.
- `/bus` - Retrieve information about specific bus routes.

### Examples

#### Home Endpoint

To check if the API is up and running, you can make a GET request to the home endpoint:

```http
GET http://localhost:8080/
```

#### Retrieve All Place data

To retrieve information about all journey place data, you can make a GET request to the `/all` endpoint:

```http
GET http://localhost:8080/all
```

You can also filter the results by specifying the `from` and `to` query parameters to narrow down your search.

#### Retrieve Specific Bus Routes

To retrieve information about specific bus routes, you can make a GET request to the `/bus` endpoint and provide the necessary query parameters, including `fromPlaceName`, `startPlaceId`, `toPlaceName`, `endPlaceId`, and `journeyDate`.

```http
GET http://localhost:8080/bus?fromPlaceName=BENGALURU&startPlaceId=1467467616730&toPlaceName=MANGALURU&endPlaceId=1467464668557&journeyDate=26/10/2023
```

## Caching

The application utilizes a caching mechanism to improve performance. Place data is cached for a specific duration (default: 24 hours) to reduce the need for frequent data retrieval.

## Contributing

Contributions to the `ksrtc-api` project are welcome! If you find any issues or have suggestions for improvements, please create a GitHub issue or submit a pull request.

## License

This project is licensed under the Apache 2.0 License. See the [LICENSE](LICENSE) file for details.

You can access the publicly hosted API at https://apiksrtc.rabil.me/.
