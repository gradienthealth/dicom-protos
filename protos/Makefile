.PHONY: build
build:
	protoc -I ./ --proto_path=./ Attributes.proto --go_out=plugins=grpc,paths=source_relative:.	
	sed -i'.bak' "s|gradienthealth|gradienthealth/$(SERVICE)|g" *.pb.go
	sed -i'.bak' "s|dicom-protos/|protos/dicom-protos/|g" *.pb.go

