run:
	kubectl create -f ./deployments/swagger-configmap.yml
	kubectl create -f ./deployments/swagger-deployment.yml
	kubectl create -f ./deployments/swagger-service.yml

stop:
	kubectl delete -f ./deployments/swagger-configmap.yml
	kubectl delete -f ./deployments/swagger-deployment.yml
	kubectl delete -f ./deployments/swagger-service.yml