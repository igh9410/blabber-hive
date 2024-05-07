# Makefile`

# Run the application
run:
	@echo "Running the Kubernetes cluster..."
	@kubectl apply -f ./k8s/services/
	@kubectl apply -f ./k8s/deployments/
	@kubectl apply -f ./k8s/configmaps/

shutdown:
	@echo "Shutting down the Kubernetes cluster..."

	@kubectl scale --replicas=0 deployment blabber-hive
	@kubectl scale --replicas=0 deployment broker
	@kubectl scale --replicas=0 deployment fastapi
	@kubectl scale --replicas=0 deployment grafana
	@kubectl scale --replicas=0 deployment kafka-setup
	@kubectl scale --replicas=0 deployment nginx
	@kubectl scale --replicas=0 deployment postgres
	@kubectl scale --replicas=0 deployment prometheus
	@kubectl scale --replicas=0 deployment redis
	@kubectl scale --replicas=0 deployment zookeeper