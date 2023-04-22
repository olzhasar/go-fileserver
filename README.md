# Simple File Server

A simple file server for uploading and downloading files, built using Go as a practice project to learn Golang.

## Features

- Upload a file to the server via an HTTP POST request
- Download a file from the server via an HTTP GET request

## Usage

### Run the server

```bash
go run main.go
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

## TODO

- Unit tests
