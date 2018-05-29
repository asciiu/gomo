build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/asciiu/gomo/device-service proto/device/device.proto
	find . -name 'device.pb.go' | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'
	docker build -t us.gcr.io/fomo-204904/device-service:latest .
	docker push us.gcr.io/fomo-204904/device-service:latest 

deploy:
	sed "s/{{ UPDATED_AT }}/$(shell date)/g" ./deployments/deployment.tmpl > ./deployments/deployment.yml
	kubectl delete -f ./deployments/deployment.yml
	kubectl create -f ./deployments/deployment.yml