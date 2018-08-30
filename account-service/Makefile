build:
	dep ensure -update github.com/asciiu/gomo
	protoc --proto_path=$(GOPATH)/src --go_out=plugins=micro:$(GOPATH)/src github.com/asciiu/gomo/account-service/proto/balance/balance.proto
	protoc --proto_path=$(GOPATH)/src --go_out=plugins=micro:$(GOPATH)/src github.com/asciiu/gomo/account-service/proto/account/account.proto
	docker build -t us.gcr.io/fomo-204904/account-service:dev .

update:
	dep ensure -update github.com/asciiu/gomo

stage:
	docker push us.gcr.io/fomo-204904/account-service:dev 

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
