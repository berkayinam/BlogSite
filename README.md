Medium.com dan sıkıldıgım için kendi blog sitemi yapıyorum

minikube ortamına gidiş: eval $(minikube docker-env)

:dönüş eval $(minikube docker-env -u)


kubectl apply -f auth-deployment.yaml   kubectl delete deployment auth-deployment 


çıkılan konteynerleri temizleme docker rm $(docker ps -a -q -f "status=exited")



curl -X POST http://192.168.49.2:port/register \
  -H "Content-Type: application/json" \
  -d '{"username":"berkay", "password":"sifrem123"}'

curl -X POST http://192.168.49.2:31488/login \
    -H "Content-Type: application/json" \
    -d '{"username":"berkay", "password":"sifrem123"}'

curl -X POST http://192.168.49.2:31660/posts   -H "Content-Type: application/json"   -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDUzMjc5ODQsInVzZXJuYW1lIjoieWlnaXQifQ.tfzevG5_PznKWepbJMTTS9Ju1FJPdzN_3ndmMw7uX7s"   -d '{"title": "İlk Yazı", "content": "Selam dünya!"}'
{"message":"Post created"}

Public Endpoints:
GET /teams - List all teams
GET /teams/:id - Get team details

Protected Endpoints (requires JWT token):
POST /teams - Create a new team
POST /teams/invite - Invite a user to team
POST /teams/invite/respond - Accept/reject team invitation
POST /teams/join/request - Request to join a team


curl -X POST http://192.168.49.2:31488/login -H "Content-Type: application/json" -d '{"username":"berkay", "password":"sifrem123"}'

curl -X POST http://192.168.49.2:31082/teams -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDU0MDgzMTgsInVzZXJuYW1lIjoiYmVya2F5In0.3ayFIPYbngwR9sjZbEERdlEs9d9opc0I-EeRtAA1CHA" -d '{"name": "Engineering Team", "description": "Our awesome engineering team"}'

curl -X GET http://192.168.49.2:31082/teams
[{"id":1,"name":"Engineering Team","description":"Our awesome engineering team","created_by":"berkay","creat
ed_at":"2025-04-20T11:38:44.002063Z","updated_at":"2025-04-20T11:38:44.002064Z","deleted_at":null,"members":
[{"id":1,"team_id":1,"username":"berkay","role":"admin","joined_at":"2025-04-20T11:38:44.007463Z","deleted_a
t":null}]}]


curl -X POST http://192.168.49.2:31082/teams/invite -H "Content-Type: application/json" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDU0MDgzMTgsInVzZXJuYW1lIjoiYmVya2F5In0.3ayFIPYbngwR9sjZbEERdlEs9d9opc0I-EeRtAA1CHA" -d '{"teamID": 1, "username": "john.doe"}'

'{"teamID": 1, "username": "john.doe"}'
{"error":"only team admins can send invites"}%                                                              
localhost%

✅ Authentication is working (tokens are being validated)
✅ Team creation is working
✅ Team listing is working
✅ Team membership is working (you were automatically added as admin)
✅ Authorization is working (the invite endpoint correctly checks for admin role)