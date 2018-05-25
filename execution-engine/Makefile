build:
	docker build -t us.gcr.io/fomo-204904/execution-engine:latest .
	docker push us.gcr.io/fomo-204904/execution-engine:latest 

deploy:
	sed "s/{{ UPDATED_AT }}/$(shell date)/g" ./deployments/deployment.tmpl > ./deployments/deployment.yml
	kubectl replace -f ./deployments/deployment.yml