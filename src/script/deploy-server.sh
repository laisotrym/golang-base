ROOT_DIR="/opt/safeweb/server/"

echo "--start"
echo "make executable file: "$1
chmod +x "${ROOT_DIR}"$1

echo "create soft link"
ln -sf "${ROOT_DIR}"$1 "${ROOT_DIR}safe-server"

echo "restart server"
systemctl restart sfserver.service

echo "--finish"
