#!/bin/bash
KEY="/c/Users/phanh/.ssh/id_rsa_safeweb"
USER="root"
IP="45.76.182.47"

LOCAL_DIR="/c/Projects/Safe-Server/src/bin/"
REMOTE_DIR="/opt/safeweb/"

scp -i $KEY $LOCAL_DIR$1 $USER@$IP:$REMOTE_DIR"server"
ssh -i $KEY $USER@$IP $REMOTE_DIR"script/deploy-server.sh" $1