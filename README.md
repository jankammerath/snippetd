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

| Language     | MIME Types                                                                                                                  | Container                                        |
|--------------|-----------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------|
| bash         | `application/x-sh`<br />`application/x-bash`<br />`text/x-sh`<br />`text/x-shellscript`                                     | `debian:latest`                                  |
| php          | `application/x-httpd-php`<br />`application/x-php`<br />`text/x-php`                                                        | `php:latest`                                     |
| python       | `application/x-python-code`<br />`application/x-python`,<br />`text/x-python`                                               | `python:latest`                                  |
| ruby         | `application/x-ruby`<br />`text/x-ruby`                                                                                     | `ruby:latest`                                    |
| javascript   | `application/javascript`<br />`text/javascript`<br />`application/x-javascript`                                             | `node:latest`                                    |
| go           | `application/x-go`<br />`text/x-go`<br />`text/x-go-source`                                                                 | `golang:latest`                                  |
| c            | `text/x-c`<br />`text/x-c-header`<br />`application/x-c`<br />`application/x-c-header`                                      | `gcc:latest`                                     |
| cpp          | `text/x-c++`<br />`text/x-c++-header`<br />`application/x-c++`<br />`application/x-c++-header`                              | `gcc:latest`                                     |
| csharp       | `application/x-csharp`<br />`text/x-csharp`<br />`text/x-csharp-source`                                                     | `mcr.microsoft.com/dotnet/sdk:latest`            |
| vbnet        | `text/x-vb`<br />`application/x-vb`                                                                                         | `mcr.microsoft.com/dotnet/sdk:latest`            |
| java         | `text/x-java-source`<br />`text/x-java`<br /> `application/x-java-source`<br />`application/x-java`<br />`application/java` | `openjdk:latest`                                 |
| rust         | `text/x-rust`                                                                                                               | `rust:latest`                                    |
| swift        | `text/x-swift`                                                                                                              | `swift:latest`                                   |
| typescript   | `application/typescript`<br />`text/typescript`                                                                             | `mcr.microsoft.com/devcontainers/typescript-node` |