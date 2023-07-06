# go-fileserver

A tiny file sharing server implemented in Go Programming Language

## Features

- Upload files to the server using HTTP POST request
- Download previously uploaded files using unique URLs
- Filesystem and in-memory storage backends
- Logging all incoming requests

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

The server will return a unique URL that can be used to download the file later.

### Download a file

To download a file, use the download link returned by the `/upload` endpoint:

```bash
curl -O -J -L http://localhost:8080/download/?token=gmjaeohnmbggokap
```

## Roadmap

- Redis registry
- S3 storage back-end
- Server configurations via command line, .yaml file
- File logger
