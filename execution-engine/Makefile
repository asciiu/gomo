build:
	dep ensure -update github.com/asciiu/gomo
	protoc --proto_path=$(GOPATH)/src --go_out=plugins=micro:$(GOPATH)/src github.com/asciiu/gomo/execution-engine/proto/engine/engine.proto
	find . -name 'engine.pb.go' | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'
	docker build --no-cache -t us.gcr.io/fomo-204904/execution-engine:dev -f ../execution-engine/Dockerfile ..  

update:
	dep ensure -update github.com/asciiu/gomo

stage:
	docker push us.gcr.io/fomo-204904/execution-engine:dev

run:
	sed "s/{{ UPDATED_AT }}/$(shell date)/g" ./deployments/deployment.tmpl > ./deployments/deployment.yml
	kubectl create -f ./deployments/deployment.yml
	kubectl create -f ./deployments/service.yml

stop:
	kubectl delete -f ./deployments/deployment.yml
	kubectl delete -f ./deployments/service.yml

restart: 
	kubectl delete -f ./deployments/deployment.yml
	kubectl delete -f ./deployments/service.yml
	sed "s/{{ UPDATED_AT }}/$(shell date)/g" ./deployments/deployment.tmpl > ./deployments/deployment.yml
	kubectl create -f ./deployments/deployment.yml
	kubectl create -f ./deployments/service.yml