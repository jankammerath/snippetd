# snippetd

A Linux daemon serving an API that allows executing, interpreting or compiling source code from various programming languages using containerd and containers from the docker hub.

![Hello World in PHP using Postman](doc/postman-php.png)

## Usage

The API is very simple as it only provides a basic banner in the root `/`, a list of all supported programming language runtimes under `/runtimes` and the execution of code under `/execute`.

```bash
curl -X POST \
    -H "Content-Type: text/x-php" \
    -d "<?php echo 'hello world';" \
    http://192.168.1.3:8080/execute
```

When posting a source code to the endpoint, it will check if the `Content-Type` _(MIME Type)_ is supported. If the language is supported, it'll create a temporary folder on the host, a container for the language with the temporary folder that includes the default source file name and the execution shell script. The shell scripts can be found in the [config/runtime](config/runtime). 

## Supported languages

The following programming languages are currently supported with the respective MIME types and containers.

| Language     | MIME Types                                                                                                 | Container                                        |
|--------------|------------------------------------------------------------------------------------------------------------|--------------------------------------------------|
| bash         | `application/x-sh`,<br /> `application/x-bash`,<br /> `text/x-sh`,<br /> `text/x-shellscript`                                | `debian:latest`                                  |
| php          | `application/x-httpd-php`, `application/x-php`, `text/x-php`                                               | `php:latest`                                     |
| python       | `application/x-python-code`, `application/x-python`, `text/x-python`                                       | `python:latest`                                  |
| ruby         | `application/x-ruby`, `text/x-ruby`                                                                        | `ruby:latest`                                    |
| javascript   | `application/javascript`, `text/javascript`, `application/x-javascript`                                    | `node:latest`                                    |
| go           | `application/x-go`, `text/x-go`, `text/x-go-source`                                                        | `golang:latest`                                  |
| c            | `text/x-c`, `text/x-c-header`, `application/x-c`, `application/x-c-header`                                 | `gcc:latest`                                     |
| cpp          | `text/x-c++`, `text/x-c++-header`, `application/x-c++`, `application/x-c++-header`                         | `gcc:latest`                                     |
| csharp       | `application/x-csharp`, `text/x-csharp`, `text/x-csharp-source`                                            | `mcr.microsoft.com/dotnet/sdk:latest`            |
| vbnet        | `text/x-vb`, `application/x-vb`                                                                            | `mcr.microsoft.com/dotnet/sdk:latest`            |
| java         | `text/x-java-source`, `text/x-java`, `application/x-java-source`, `application/x-java`, `application/java` | `openjdk:latest`                                 |
| rust         | `text/x-rust`                                                                                              | `rust:latest`                                    |
| swift        | `text/x-swift`                                                                                             | `swift:latest`                                   |
| typescript   | `application/typescript`, `text/typescript`                                                                | `mcr.microsoft.com/devcontainers/typescript-node` |