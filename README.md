# Context Free Grammar API

> Inspired by CS 420: Structure of Programming Languages from George Fox University assignment to implement context free grammar generation.

## Overview

The Context Free Grammar API provides a service inspired by the assignment from CS 420: Structure of Programming Languages at George Fox University. The goal is to implement a context-free grammar generation mechanism using a grammar file that defines grammar rules.

## Features

- Connects to a MongoDB instance to store and retrieve grammar configurations.
- Exposes an HTTP API for generating strings based on defined grammar rules.
- Supports flexible and customizable grammar configurations.
- Uses Vercel for seamless deployment and scaling.

## Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/yourusername/go.resumes.guide.git
    cd go.resumes.guide
    ```

2. **Configure the environment:**

    Copy or rename `.env.template` to `.env` and adjust the appropriate settings like MongoDB URI and server address.

3. **Install dependencies:**

    Ensure you have Go installed (minimum version 1.19) and use `go mod` to install dependencies:

    ```bash
    go mod tidy
    ```

## Usage

### Seeding the Database

Before running the API, seed the database with initial grammar configurations:

```bash
cd cmd/seed
go run main.go
```

The server will start on the address specified in your .env file (e.g., :8080).

## API Endpoints
- Generate Grammar-Based Text
- Endpoint: `/`
- Method: `GET`
- Response: JSON containing the generated text based on the grammar.


### Example Request
```
curl -X GET http://localhost:8080
```

### Example Response
```
{
  "message": "Generated text based on the grammar",
  "status": "OK"
}
```

## Deployment
The application is configurable to deploy on Vercel. Use the `vercel.json` configuration file to adjust deployment settings. Ensure the necessary build and route settings are correctly defined in vercel.json.

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request with any enhancements or bug fixes.


## Contact
For issues, improvements, or inquiries, create an issue in the GitHub repository or contact the maintainers.
