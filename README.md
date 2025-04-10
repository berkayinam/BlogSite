Medium.com dan sıkıldıgım için kendi blog sitemi yapıyorum

minikube ortamına gidiş: eval $(minikube docker-env)

:dönüş eval $(minikube docker-env -u)

ssh -L 8080:192.168.49.2:31488 binam@45.9.30.161

kubectl apply -f auth-deployment.yaml   kubectl delete deployment auth-deployment 

çıkılan konteynerleri temizleme : docker rm $(docker ps -a -q -f "status=exited")



curl -X POST http://192.168.49.2:port/register \
  -H "Content-Type: application/json" \
  -d '{"username":"berkay", "password":"sifrem123"}'


curl -X POST http://192.168.49.2:31488/login \
    -H "Content-Type: application/json" \
    -d '{"username":"berkay", "password":"sifrem123"}'
    
{"token":"eyJh***********7GebnjYe4BJ28-4_LiAgA"}

register ve  login successss. yesssssssss


curl -X POST http://192.168.49.2:31660/posts   -H "Content-Type: application/json"   -H "Authorization: Bearer eyJhbG**********vVlJg"   -d '{"title": "İlk Yazı", "content": "Selam dünya!"}'
{"message":"Post created"}

yesss post created

![image](https://github.com/user-attachments/assets/c1639a1e-8e33-4fce-8301-cbe14667ff90)
