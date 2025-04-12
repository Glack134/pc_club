# Запрашиваем IP сервера у пользователя
$serverIP = Read-Host "192.168.1.14:50051"

# Создаем минимальный конфиг
@"
server_ip: "$serverIP"
pc_name: "$env:COMPUTERNAME"
"@ | Out-File -FilePath "C:\Program Files\PC_Club_Client\config.yaml"

# Запускаем клиент
Start-Process -FilePath "C:\Program Files\PC_Club_Client\pc_club_client.exe"