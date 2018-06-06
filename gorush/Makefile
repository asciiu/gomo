start:
	kubectl create -f ./deployments/gorush-configmap.yml
	kubectl create -f ./deployments/gorush-deployment.yml
	kubectl create -f ./deployments/gorush-service.yml

stop:
	kubectl delete -f ./deployments/gorush-configmap.yml
	kubectl delete -f ./deployments/gorush-deployment.yml
	kubectl delete -f ./deployments/gorush-service.yml