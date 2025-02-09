# Go Client-Server Application

This project contains a simple client and server written in Go that communicate using a binary protocol with `gob` encoding. The server challenges the client with a proof-of-work (POW) task and, upon successful validation, sends a random quote back to the client.

## Components

### Server

- **Purpose:**
    - Listens on TCP port 8080.
    - When a client connects, sends a POW challenge.
    - Waits for the client to respond with a solution.
    - Verifies the solution using a built-in proof-of-work mechanism.
    - If the solution is correct, returns a random quote.

- **Key Packages:**
    - `encoding/gob` for binary serialization.
    - Custom packages (`pow`, `protocol`, and `quotes`) for generating challenges, handling message types, and retrieving quotes.

### Client

- **Purpose:**
    - Connects to the server at a specified address (by default `localhost:8080`).
    - Receives the challenge from the server.
    - Computes the POW solution.
    - Sends the solution back to the server.
    - Receives and displays the quote if the solution is accepted.

- **Configuration:**
    - Uses a command-line flag `-server` to set the server address.
    - The client Dockerfile is set up to accept the server address as a build argument (which is then passed as an environment variable).

## Building and Running

This project uses multi-stage Dockerfiles for both the client and server to produce minimal runtime images.

### Using Docker

#### Server

1. **Build the Server Image:**  
   Navigate to the `server` directory and run:
   ```bash
   docker build -t myserver .
   ```

2. **Run the Server Container:**
     
    Start the `server` and map port `8080`:
   ```bash
   docker run -p 8080:8080 myserver
   ```
    The server will now listen on port 8080.

#### Client

1. **Build the Client Image:**
   
    Navigate to the `client` directory and run:
    ```bash
    docker build -t myclient --build-arg SERVER_ADDR=server:8080 .
    ```
    This passes the server address as a build argument, which is set in the container as the environment variable `SERVER_ADDR`.
2. **Run the Client Container:**
   
   Start the client:

    ```bash
   docker run myclient
    ```
   The client binary will be executed with the command-line flag `-server server:8080` (using the value from the build argument).

### Using Docker Compose

You can also launch both services together using Docker Compose. Below is an example `docker-compose.yml` file:

```yaml
version: "3.8"

services:
  server:
    build:
      context: ./server
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    networks:
      - app-network

  client:
    build:
      context: ./client
      dockerfile: Dockerfile
      args:
        SERVER_ADDR: server:8080
    depends_on:
      - server
    networks:
      - app-network

networks:
  app-network:
    driver: bridge
```

To build and run both containers, execute:

```bash
docker-compose up --build
```
This setup ensures that the client connects to the server by its service name (`server`) on port `8080`.

## Source Code Overview
- **Server Code**:

    Located in the `./cmd/server/main.go` directory, the server listens for incoming TCP connections, sends a POW challenge, verifies the solution, and sends a random quote on success.

- **Client Code**:

    Located in the `./cmd/client/main.go` file, the client accepts a server address via the `-server` flag, connects to the server, processes the challenge, and handles the quote response.

## Customization
- **Server Address**:

    Modify the `SERVER_ADDR` build argument in the client Dockerfile or override it in your Docker Compose file if your server runs on a different address or port.

- **Ports**:

    Adjust the port mappings in the Docker run commands or Docker Compose file to suit your environment.

## License
This project is licensed under the MIT License.

```yaml
---
This `README.md` provides an overview of the project's purpose, details each component's role, and explains how to build and run the client and server using both Docker and Docker Compose.
```