#!/bin/bash
# Скрипт управления доступом к игровым ПК
# Использование: ./pc_control.sh [lock|unlock] <pc-id>

SERVER_ADDR="192.168.1.14:50051"
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDM0NDEzODAsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjoiYWRtaW4ifQ.yMFyqjXqCUWt0oA9nJ7o0tZbq5t3YUaMOo3Mx6qydZU"  # Замените на реальный токен

# Проверка аргументов
if [ $# -ne 2 ]; then
    echo "Ошибка: Неверное количество аргументов"
    echo "Правильное использование: $0 {lock|unlock} pc-id"
    exit 1
fi

ACTION=$1
PC_ID=$2

case $ACTION in
    lock)
        echo "Блокировка ПК $PC_ID..."
        grpcurl -plaintext -d "{\"pc_id\":\"$PC_ID\", \"auth_token\":\"$TOKEN\"}" \
            $SERVER_ADDR rpc.AdminService/RevokeAccess
        ;;
    unlock)
        echo "Разблокировка ПК $PC_ID..."
        grpcurl -plaintext -d "{\"user_id\":\"temp-user\", \"pc_id\":\"$PC_ID\", \"minutes\":60, \"auth_token\":\"$TOKEN\"}" \
            $SERVER_ADDR rpc.AdminService/GrantAccess
        ;;
    *)
        echo "Ошибка: Неизвестное действие '$ACTION'"
        echo "Допустимые действия: lock, unlock"
        exit 1
        ;;
esac

# Проверка результата
if [ $? -eq 0 ]; then
    echo "Операция выполнена успешно"
else
    echo "Ошибка при выполнении операции"
    exit 1
fi