FROM golang:alpine
WORKDIR /go/src/github.com/gradienthealth/dicom-protos
COPY . .
RUN apk add --no-cache make protobuf git
RUN go get -u github.com/golang/protobuf/protoc-gen-go
RUN mkdir protos
RUN make
CMD ./dicom-protos