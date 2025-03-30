#!/bin/bash

# Генерация секретного ключа (если еще нет)
if [ ! -f .secret ]; then
  openssl rand -base64 32 > .secret
fi

SECRET=$(cat .secret)

# Генерация токена
go run cmd/token_gen/main.go -secret "$SECRET" > .token

# Создание config.ini
cat > bin/config.ini <<EOL
[client]
server_address = 192.168.1.14:50051  # Замените на IP сервера
pc_id = ${1:-windows-pc}             # Имя ПК (можно передать как аргумент)
auth_token = $(cat .token)           # Автоматически подставляем токен
EOL

echo "Клиент готов в папке bin/"
ls -lh bin/client.exe bin/config.ini