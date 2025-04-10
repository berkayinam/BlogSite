Medium.com dan sıkıldıgım için kendi blog sitemi yapıyorum

minikube ortamına gidiş: eval $(minikube docker-env)

:dönüş eval $(minikube docker-env -u)


kubectl apply -f auth-deployment.yaml   kubectl delete deployment auth-deployment 


çıkılan konteynerleri temizleme docker rm $(docker ps -a -q -f "status=exited")



curl -X POST http://192.168.49.2:32669/reg
ister \
  -H "Content-Type: application/json" \
  -d '{"username":"berkay", "password":"sifrem123"}'


curl -X POST http://192.168.49.2:32669/login   -H "Content-Type: application/json"   -d '{"username":"berkay", "password":"sifrem
123"}'
{"token":"eyJh***********7GebnjYe4BJ28-4_LiAgA"}



register ve  login successss. yesssssssss