# BlogSite

Medium.com benzeri gelişmiş bir blog platformu. Go backend, Kubernetes ve modern frontend teknolojileri kullanılarak geliştirilecektir.

## 🚀 Proje Hakkında

Bu proje, modern web teknolojilerini kullanarak kapsamlı bir blog platformu oluşturmayı ve DevOps pratiklerini uygulamayı amaçlamaktadır. Proje, backend'de Go dilini, veritabanı olarak MongoDB'yi ve container orchestration için Kubernetes'i kullanmaktadır.

## 💻 Teknolojiler

### Backend
- Go (Golang)
- MongoDB
- RESTful API

### DevOps & Deployment
- Kubernetes
- Docker
- Container Orchestration

### Veritabanı
- MongoDB

## 🛠️ Kurulum

### Gereksinimler
- Go 1.16 veya üzeri
- MongoDB
- Docker (opsiyonel)
- Kubernetes (opsiyonel)

### Yerel Geliştirme Ortamı
1. Projeyi klonlayın:
```bash
git clone https://github.com/[kullanıcı-adı]/BlogSite.git
cd BlogSite
```

2. Bağımlılıkları yükleyin:
```bash
go mod download
```

3. MongoDB bağlantısını yapılandırın:
- MongoDB'yi yerel olarak çalıştırın veya MongoDB Atlas kullanın
- Bağlantı bilgilerini yapılandırma dosyasında güncelleyin

4. Uygulamayı çalıştırın:
```bash
go run main.go
```

## 📝 Özellikler

- [ ] Kullanıcı Yönetimi
  - [ ] Kayıt
  - [ ] Giriş/Çıkış
  - [ ] Profil Yönetimi

- [ ] Blog Yazıları
  - [ ] Yazı Oluşturma
  - [ ] Düzenleme
  - [ ] Silme
  - [ ] Görüntüleme

- [ ] Kategoriler ve Etiketler
  - [ ] Kategori Bazlı Filtreleme
  - [ ] Etiket Bazlı Arama

- [ ] Etkileşim
  - [ ] Beğeni
  - [ ] Yorum
  - [ ] Paylaşım

## 🔄 API Endpoints

### Blog Yazıları
- `GET /api/posts` - Tüm yazıları listele
- `GET /api/posts/{id}` - Belirli bir yazıyı getir
- `POST /api/posts` - Yeni yazı oluştur
- `PUT /api/posts/{id}` - Yazı güncelle
- `DELETE /api/posts/{id}` - Yazı sil

### Kullanıcılar
- `POST /api/users/register` - Yeni kullanıcı kaydı
- `POST /api/users/login` - Kullanıcı girişi
- `GET /api/users/profile` - Kullanıcı profili
- `PUT /api/users/profile` - Profil güncelleme

## 🔐 Güvenlik

- JWT tabanlı kimlik doğrulama
- Şifrelenmiş kullanıcı bilgileri
- Rate limiting
- CORS politikaları

## 🚀 Deployment

### Docker ile Deployment
```bash
# Docker image oluştur
docker build -t blogsite .

# Container'ı çalıştır
docker run -p 8080:8080 blogsite
```

### Kubernetes ile Deployment
```bash
# Kubernetes deployment'ı uygula
kubectl apply -f k8s/deployment.yaml

# Service'i oluştur
kubectl apply -f k8s/service.yaml
```

## 📈 Gelecek Özellikler

- [ ] Gelişmiş editör desteği
- [ ] Çoklu dil desteği
- [ ] Elasticsearch entegrasyonu
- [ ] Otomatik backup sistemi
- [ ] Analytics dashboard
- [ ] Email bildirim sistemi

## 🤝 Katkıda Bulunma

1. Bu repository'yi fork edin
2. Feature branch'i oluşturun (`git checkout -b feature/AmazingFeature`)
3. Değişikliklerinizi commit edin (`git commit -m 'Add some AmazingFeature'`)
4. Branch'inizi push edin (`git push origin feature/AmazingFeature`)
5. Pull Request oluşturun

## 📝 Lisans

Bu proje [MIT](LICENSE) lisansı altında lisanslanmıştır.

## 👨‍💻 Geliştirici

[Geliştirici Adı] - [GitHub Profili]

---

⭐️ Bu projeyi beğendiyseniz yıldız vermeyi unutmayın!



---- Notlar ----
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

curl -X POST http://192.168.49.2:31660/posts   -H "Content-Type: application/json"   -H "Authorization: Bearer eyJhbG**********vVlJg"   -d '{"title": "İlk Yazı", "content": "Selam dünya!"}'
{"message":"Post created"}

![image](https://github.com/user-attachments/assets/c1639a1e-8e33-4fce-8301-cbe14667ff90)
