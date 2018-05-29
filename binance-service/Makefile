build:
	docker build -t us.gcr.io/fomo-204904/binance-service:latest .
	docker push us.gcr.io/fomo-204904/binance-service:latest 

deploy:
	sed "s/{{ UPDATED_AT }}/$(shell date)/g" ./deployments/deployment.tmpl > ./deployments/deployment.yml
	kubectl create -f ./deployments/deployment.yml
	kubectl create -f ./deployments/service.yml

stop:
	kubectl delete -f ./deployments/deployment.yml
	kubectl delete -f ./deployments/service.yml