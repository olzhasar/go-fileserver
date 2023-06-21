# go-fileserver

A tiny file server implemented in Go Programming Language

## Features

- Upload a file to the server using HTTP POST request
- Download a file from the server using HTTP GET request
- Log every incoming request with duration info to `stdout`

## Usage

### Run the server

```bash
go run .
```

The server will start on port 8080.

### Upload a file

To upload a file, send a POST request to `/upload` with a form-data containing the file:

```bash
curl -X POST -F "file=@/path/to/your/file.txt" http://localhost:8080/upload
```

Replace `/path/to/your/file.txt` with the path to the file you want to upload.

### Download a file

To download a file, send a GET request to `/download?filename={filename}`:

```bash
curl -O -J -L http://localhost:8080/download/?filename=example.txt
```

Replace `filename.txt` with the name of the file you want to download.

## Roadmap

- Unique urls for each upload persisted in a database (e.g. SQLite)
- Multiple options for persisting uploads (SQLite, PostgreSQL, Redis)
- S3 storage back-end
- File logger
- Server configurations via command line, .yaml file
