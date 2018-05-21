build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/asciiu/gomo/device-service proto/device/device.proto
	find . -name 'device.pb.go' | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'