kubectl delete -f kubernetes/

# Comment Service için
cd comment-service/
docker build -t comment-service:latest .

# Team Service için
cd ../team-service
docker build -t team-service:latest .

# auth Service için
cd ../auth-service
docker build -t auth-service:latest .

# post Service için
cd ../post-service
docker build -t post-service:latest .

cd ../kubernetes
# MongoDB ve diğer bağımlılıkları başlat
kubectl apply -f postgres-secret.yaml
kubectl apply -f postgres-pv-pvc.yaml
kubectl apply -f postgres-deployment.yaml
kubectl apply -f postgres-service.yaml

# Auth service'i başlat
kubectl apply -f auth-deployment.yaml
kubectl apply -f auth-service.yaml

# Team service'i başlat
kubectl apply -f team-deployment.yaml
kubectl apply -f team-service.yaml

# Comment service'i başlat
kubectl apply -f comment-deployment.yaml
kubectl apply -f comment-service.yaml

#post servici baslat
kubectl apply -f post-deployment.yaml
kubectl apply -f post-service.yaml
