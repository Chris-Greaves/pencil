<div align="center">
    <img src="README.md.files/pencil2.jpg" alt="Pencil logo" height=200>
    <h1>Pencil</h1>
    <img src="https://img.shields.io/github/go-mod/go-version/chris-greaves/pencil?filename=go.mod">
    <a href="https://github.com/Chris-Greaves/pencil/actions/workflows/go.yml"><img src="https://github.com/chris-greaves/pencil/actions/workflows/go.yml/badge.svg" alt="Build Workflow status badge"></a>
    <a href="https://github.com/Chris-Greaves/pencil/releases"><img src="https://github.com/chris-greaves/pencil/actions/workflows/release.yml/badge.svg" alt="Release Workflow status badge"></a>
    <a href="https://github.com/Chris-Greaves/pencil/releases"><img src="https://img.shields.io/github/v/release/chris-greaves/pencil?label=Latest%20Release" alt="Latest Release badge"></a>
    <a href="https://github.com/Chris-Greaves/pencil/pkgs/container/pencil"><img src="https://github.com/chris-greaves/pencil/actions/workflows/docker-publish.yml/badge.svg" alt="Docker Workflow status badge"></a>
</div>

# 

A simple find and replace tool using Go's `text/template` library, useful for filling out config files and scripts, especially when paired with docker-compose projects.

## Getting started

Pencil is available as a docker image or executable.

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

## How do Variables Work?

In order for a Template engine to be useful, you've gotta give it some values! Pencil handles variables in a few ways; below is a table showing how to input variables and then use them.

| Type | Input | Usage |
| --- | --- | --- |
| Direct | `pencil -v Name=Value` | `{{ var "Name" }}` or `{{ .Var.Name }}` |
| Environment | `Name=Value pencil` | `{{ env "Name" }}` or `{{ .Env.Name }}` |

> NOTE: More variable types are likely to come in the future.

## In-place or Template?

**In-place**  
By this, we mean a file whose content is treated as a template, Parsed and Executed, and then written back into the file. Because the file's content is overwritten, this is a permanent change and cannot be re-run with different variables. This strategy is desirable in specific docker-compose scenarios where you don't want to specify the variables every time you run a restart or rebuild. A brief demonstration can be found below:

```bash
echo "Hello {{ var "Name" }}!" > test.txt
pencil -v Name=Chris test.txt
cat test.txt # Outputs: Hello Chris!
pencil -v Name=Kyle test.txt
cat test.txt # Outputs: Hello Chris!
```

Because the file no longer contains the template syntax, it is unchanged when run a second time with a different value.

**Template**  
This is when a template file is used and not changed, even after Pencil has executed it. So, instead of the file's content being overwritten, the resulting content is written into a new file. This strategy is useful when you often want to make config changes and don't mind repeatedly supplying the variables. In docker-compose situations, the generated file is recreated every time the Pencil service is re-run. A brief demonstration can be found below:

```bash
echo "Hello {{ var "Name" }}!" > test.txt.gotmpl
ls                  # Outputs: test.txt.gotmpl
pencil -v Name=Chris test.txt.gotmpl
ls                  # Outputs: test.txt, test.txt.gotmpl
cat test.txt        # Outputs: Hello Chris!
cat test.txt.gotmpl # Outputs: Hello {{ var "Name" }}!
pencil -v Name=Kyle test.txt
ls                  # Outputs: test.txt, test.txt.gotmpl
cat test.txt        # Outputs: Hello Kyle!
cat test.txt.gotmpl # Outputs: Hello {{ var "Name" }}!
```

As you can see, a file was created alongside the template with the resulting content from executing it. The file name is the same as the template without the `.gotmpl` extension. The template file is unchanged, allowing you to re-run Pencil to update the generated file with the new variable.