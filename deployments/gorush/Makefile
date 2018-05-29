start:
	kubectl create -f ./gorush-configmap.yml
	kubectl create -f ./gorush-deployment.yml
	kubectl create -f ./gorush-service.yml

stop:
	kubectl delete -f ./gorush-configmap.yml
	kubectl delete -f ./gorush-deployment.yml
	kubectl delete -f ./gorush-service.yml