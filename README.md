# Pencil

A simple find and replace tool using Go's `text/template` library, useful for filling out config files and scripts, especially when paired with docker-compose projects.

## Getting started

Pencil is available as a executable or docker image.

### Docker Compose

1. Create a file called `test.txt` containing the following text:  
  ```Hello {{ env "Name" }}! This is {{ var "AppName" }}.```
2. Create `compose.yaml` with the following content:  
   ```yaml
    services:
      pencil:
        image: ghcr.io/chris-greaves/pencil:v0.0.1
        environment:
          - Name=Chris          # Environment Variables available using {{ env "KEY" }}
        command:
          - "-v AppName=Pencil" # Direct Variables available using {{ var "KEY" }}
          - "/mnt/test.txt"     # Tell Pencil what file or folders to process.
        volumes:
          - ./test.txt:/mnt/test.txt # Mount the files or folders you want to process here.
   ```
3. Run `docker compose up`
4. Now check the contents of `test.txt` (e.g. `cat test.txt`)

### Docker

1. Create a file called `test.txt` containing the following text:  
  ```Hello {{ env "Name" }}!```
2. Run this docker command  
   Windows: `docker run -e Name=Chris -v .\test.txt:/mnt/test.txt ghcr.io/chris-greaves/pencil:latest /mnt/test.txt`  
   Linux & MacOS: `docker run -e Name=Chris -v ./test.txt:/mnt/test.txt ghcr.io/chris-greaves/pencil:latest /mnt/test.txt`
3. Now check the contents of `test.txt` (e.g. `cat test.txt`)

### Executable:

1. Download the latest executable for your OS from the [releases page](https://github.com/Chris-Greaves/pencil/releases).

2. Extract the `.exe` and place is somewhere in your PATH.
3. Create a file called `test.txt` containing the following text:  
  ```Hello {{ var "Name" }}!```
4. Run `pencil -v Name=Chris test.txt`
5. Now check the contents of `test.txt` (e.g. `cat test.txt`)

