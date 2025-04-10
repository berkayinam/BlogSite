eval $(minikube docker-env) dönüş eval $(minikube docker-env -u)
docker rm $(docker ps -a -q -f "status=exited")