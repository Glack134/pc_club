#!/bin/bash
# Установка клиента как systemd сервиса

CLIENT_BIN="/usr/local/bin/gameclient"
CONFIG_DIR="/etc/gameaccess"

mkdir -p $CONFIG_DIR
cp gameclient $CLIENT_BIN
cp config.json $CONFIG_DIR

cat > /etc/systemd/system/gameclient.service <<EOF
[Unit]
Description=Game Access Client
After=network.target

[Service]
ExecStart=$CLIENT_BIN
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOF

systemctl enable gameclient
systemctl start gameclient