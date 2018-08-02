# dicom-protos
[![Build Status](https://travis-ci.org/gradienthealth/dicom-protos.svg?branch=master)](https://travis-ci.org/gradienthealth/dicom-protos) 
[![Go Report Card](https://goreportcard.com/badge/github.com/gradienthealth/dicom-protos)](https://goreportcard.com/report/github.com/gradienthealth/dicom-protos) 
[![GoDoc Reference](https://godoc.org/github.com/gradienthealth/dicom-protos?status.svg)](https://godoc.org/github.com/gradienthealth/dicom-protos)

This repository contains generated protocol buffer representations of every DICOM attribute and module. These protocol buffers were generated using the software included in this repository and the [innolitics JSON dump](https://github.com/innolitics/dicom-standard) of the entire DICOM standard. 

Check out our [Golang DICOM parser repo](https://github.com/gradienthealth/dicom) as well, if you're interested in medical imaging. 

# Quick Start
## Generate protocol buffers
### Native
Assuming you have golang and make installed:
1. Clone this repository using `git clone --recursive` into your $GOPATH or `go get github.com/gradienthealth/dicom-protos`
2. Simply run `make run` and the protos will be deposited in the `protos` directory. 

### Docker
If you use docker you can simply do the following to regenerate the protocol buffers:
1. Clone this repository using `git clone --recursive`
2. :eyes:Optional: Update the `dicom-standard` submodule if there are new JSON definitions to pull
3. Build the container `docker build . -t gradienthealth/dicom-protos`
3. Run the command below, replacing $PWD/OUTPUT_PROTOS with the path on your local machine you would like the generated protos to be deposited in. 
   ```sh
   docker run -v $PWD/OUTPUT_PROTOS:/go/src/github.com/gradienthealth/dicom-protos/protos -it gradienthealth/dicom-protos make run
   ```

## Code generation
You can generate native representations of these DICOM protocol buffers in the [language of your choice](https://developers.google.com/protocol-buffers/docs/proto3#generating). 

For example to generate golang representations of these data structures, you can simply run:
```sh
protoc --proto_path=. --go_out=:$GOPATH/src *.proto
```
