Приложение go-video-hosting - это рефакторинг бэкенда ранее написанного видеохостинга http://www.geek-tube.ru
Для запуска необходимо проделать следующее:

- В корне приложения создать папку data.
- Запустить докер с базой данных, выполнив команду
  `sudo docker run --name=go-video-hosting -e POSTGRES_PASSWORD=<password> -p 5436:5432 -d --rm -v ./data:/var/lib/postgresql/data postgres`.
- Если развёртывание задачи происходит первый раз и база данных еще не создавалась, создать базу данных командой `sudo docker exec -it go-video-hosting psql -U postgres -c "CREATE DATABASE go_video_hosting;"`
- Запустить миграции командой `sudo migrate -path ./schema -database "postgresql://postgres:<password>@localhost:5436/go_video_hosting?sslmode=disable" -verbose up`.
- Теперь можно запускать приложение `sudo export $(cat .env | xargs) go run ./cmd/main.go`.
