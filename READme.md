jadi untuk buat inisialisasai golang nya dulu dengan perintah :

go mod init nama_folder_nya

setelah itu install fiber di dalam folder tersebut :

go get github.com/gofiber/fiber/v3


//untuk beri tag di github ketika ada update di file proto
git tag v1.0 
git push origin v1.0   


//kode untuk mendapatkan update grpc berdasarkan tag di github
go get github.com/revivalxvr/protobuf/protogen/product@v0.0.3


//untuk ngronk pake localhost 8084(bisa disesuaikan port nya)
ngrok http 8084 --url https://available-iciness-heaving.ngrok-free.dev