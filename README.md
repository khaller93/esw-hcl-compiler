# ESW HCL Compiler

This customized compiler is designed to convert HCL (HashiCorp Configuration Language) configuration files into
JSON format, but also to resolves IRIs (Internationalized Resource Identifiers) that incorporate a prefix. This
functionality is essential for the latest version of [ESW](https://github.com/khaller93/es-web-app).

## Build

Requirements:
* Go (>= 1.13)

```
go mod vendor && go build
```

### Docker Image

This repository includes a Dockerfile that enables the creation of a compact Docker image, as demonstrated by the
command below.

```
docker build . --build-arg ESW_COMP_VERSION=latest -t yyyyy/esw-hcl-compiler:latest
```

## Run 

The build process generates a binary tailored to your operating system and architecture. Afterward, you can run the
command outlined below:

```
./esw-hcl-compiler <configuration-directory> <output-directory>
```

It will read in all configuration files that can be found in the given directory, and then write the JSON result of the
compilation to the output directory. This directory can then be handed to the [ESW](https://github.com/khaller93/es-web-app).

## Contact

* Kevin Haller - [kevin.haller@tuwien.ac.at](mailto:kevin.haller@tuwien.ac.at)