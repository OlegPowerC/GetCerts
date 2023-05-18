#!/bin/sh
echo "stop httpd"
service httpd stop
echo "copy files"
cp cert.cer /etc/httpd/conf/cert.crt
cp key.key /etc/httpd/conf/key.key
echo "start httpd"
service httpd start

