#!/bin/bash

# Build the team service Docker image
cd ../team-service
docker build -t team-service:latest .

# Apply the team service deployment and service
cd ../kubernetes
kubectl apply -f team-deployment.yaml
kubectl apply -f team-service.yaml

# Wait for the deployment to be ready
echo "Waiting for Team Service deployment to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/team-service

echo "Team Service deployment complete!" 