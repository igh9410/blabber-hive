# Makefile`

# Run the application
run:
	@echo "Setting up the stateful infrastructures..."
	@echo "Running the Kubernetes cluster..."
	@kubectl apply -k ./k8s
shutdown:
	@echo "Shutting down the Kubernetes cluster..."

	@kubectl scale --replicas=0 deployment blabber-hive
	@kubectl scale --replicas=0 statefulset broker
	@kubectl scale --replicas=0 deployment fastapi
	@kubectl scale --replicas=0 deployment grafana
	@kubectl scale --replicas=0 deployment nginx
	@kubectl scale --replicas=0 statefulset postgres
	@kubectl scale --replicas=0 deployment prometheus
	@kubectl scale --replicas=0 deployment redis
	@kubectl scale --replicas=0 deployment zookeeper
	
#	@echo "Shutting down the stateful infrastructures..."
#	@docker compose -f ./docker-compose-k8s.yml down

clean:
	@echo "Deleting all Kubernetes resources..."
	@kubectl delete all --all

stop-cluster:
	@echo "Disabling Kubernetes in Docker Desktop..."
	@osascript -e 'tell application "Docker" to set kubernetes enabled to false'