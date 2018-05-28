build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/asciiu/gomo/user-service proto/user/user.proto
	docker build -t us.gcr.io/fomo-204904/user-service:latest .
	docker push us.gcr.io/fomo-204904/user-service:latest 

deploy:
	sed "s/{{ UPDATED_AT }}/$(shell date)/g" ./deployments/deployment.tmpl > ./deployments/deployment.yml
	kubectl delete -f ./deployments/deployment.yml
	kubectl create -f ./deployments/deployment.yml