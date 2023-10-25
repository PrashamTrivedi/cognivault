# My Go Project

This project is a Go API that allows users to create collections, add data points to those collections, and perform queries on the data. It uses an SQLite database to store the data.

## Project Structure

The project has the following file structure:

```
my-go-project
├── api
│   ├── handlers.go
│   ├── routes.go
│   └── swagger.go
├── collections
│   ├── collection.go
│   ├── data_point.go
│   └── tag.go
├── database
│   └── database.go
├── go.mod
├── go.sum
├── main.go
├── README.md
└── utils
    ├── file.go
    └── response.go
```

The files in the project are organized as follows:

- `api/handlers.go`: This file contains the HTTP request handlers for the API endpoints.
- `api/routes.go`: This file sets up the routes for the API endpoints using the `chi` router.
- `api/swagger.go`: This file serves the Swagger UI for the API documentation.
- `collections/collection.go`: This file contains the `Collection` struct and methods for working with collections.
- `collections/data_point.go`: This file contains the `DataPoint` struct and methods for working with data points.
- `collections/tag.go`: This file contains the `Tag` struct and methods for working with tags.
- `database/database.go`: This file contains functions for connecting to the SQLite database and executing SQL queries.
- `utils/file.go`: This file contains functions for reading files from disk.
- `utils/response.go`: This file contains functions for creating HTTP responses.

## API Endpoints

The API has the following endpoints:

- `POST /collections`: Creates a new collection.
- `POST /collections/{collectionName}/datapoints`: Adds a new data point to a collection.
- `GET /collections/{collectionName}/datapoints`: Retrieves data points from a collection based on a query.
- `PUT /collections/{collectionName}/tags/{tagName}`: Updates a tag in a collection.
- `DELETE /collections/{collectionName}/tags/{tagName}`: Deletes a tag from a collection.
- `PUT /collections/{collectionName}`: Updates a collection.
- `DELETE /collections/{collectionName}`: Deletes a collection.
- `GET /collections/{collectionName}/tags`: Retrieves tags from a collection.
- `GET /collections/{collectionName}/tags/{tagName}/datapoints`: Retrieves data points from a tag.

## Dependencies

The project uses the following dependencies:

- `github.com/go-chi/chi`: A lightweight router for Go HTTP services.
- `github.com/go-chi/cors`: Middleware for setting up CORS headers.
- `github.com/swaggo/http-swagger`: Middleware for serving the Swagger UI.
- `github.com/mattn/go-sqlite3`: A SQLite driver for Go.

## Running the Project

To run the project, you will need to have Go 1.21 installed. Clone the repository and run the following command:

```
go run main.go
```

This will start the API server on port 8080. You can then use a tool like `curl` or a web browser to interact with the API endpoints.