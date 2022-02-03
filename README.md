# ESW HCL Compiler

A compiler that transforms HCL configuration files into JSON, and resolves IRIs that 
contain a prefix. It is required for the current version ESW.

## Build

Requirements:
* Go (>= 1.13)

```
go mod vendor && go build
```

### Docker Image

This repository provides a Dockerfile with which a small Docker image can be built as shown
in the following command.

```
docker build . --build-arg ESW_COMP_VERSION=latest -t yyyyy/esw-hcl-compiler:latest
```
Built images can be found [here](https://hub.docker.com/repository/docker/khaller/esw-hcl-compiler/tags).

## Run 

The build creates a binary for your operating system and architecture. Then, you can
execute the following command 

```
./esw-hcl-compiler <configuration-directory> <output-directory>
```

It will read in all configuration files that can be found in the given directory,
and then write the JSON result of the compilation to the output directory. This directory
can then be handed to the exploratory search application.

## Contact

* Kevin Haller - [kevin.haller@tuwien.ac.at](mailto:kevin.haller@tuwien.ac.at)