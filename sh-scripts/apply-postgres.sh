#!/bin/bash

# Create the secret first
kubectl apply -f postgres-secret.yaml

# Create PV and PVC
kubectl apply -f postgres-pv-pvc.yaml

# Create the deployment
kubectl apply -f postgres-deployment.yaml

# Create the service
kubectl apply -f postgres-service.yaml

# Wait for the deployment to be ready
echo "Waiting for PostgreSQL deployment to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/postgres

echo "PostgreSQL setup complete!" 