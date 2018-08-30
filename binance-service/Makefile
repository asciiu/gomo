build:
	dep ensure -update github.com/asciiu/gomo
	protoc --proto_path=$(GOPATH)/src --go_out=plugins=micro:$(GOPATH)/src github.com/asciiu/gomo/binance-service/proto/balance/balance.proto
	protoc --proto_path=$(GOPATH)/src --go_out=plugins=micro:$(GOPATH)/src github.com/asciiu/gomo/binance-service/proto/binance/binance.proto
	find . -name 'binance.pb.go' | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'
	find . -name 'balance.pb.go' | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'
	docker build --no-cache -t us.gcr.io/fomo-204904/binance-service:dev .

update:
	dep ensure -update github.com/asciiu/gomo

stage:
	docker push us.gcr.io/fomo-204904/binance-service:dev 

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