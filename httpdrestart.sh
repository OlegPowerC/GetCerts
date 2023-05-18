#!/bin/sh
echo "stop httpd"
service httpd stop
echo "copy files"
cp yourcert.cer /etc/httpd/conf/cert.crt
cp yourkey.key /etc/httpd/conf/key.key
echo "start httpd"
service httpd start

