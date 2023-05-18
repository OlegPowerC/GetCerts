# Синхронизация сертификатов WAF и сервера
License - GPL
### Как это работает
WAF копирует по SFTP сертификаты на промежуточный сервер, в папку доступную по HTTPS
На сервере по cron запускается утилита, которая в случае изменения даты файла сертификата, копирует сертификат, закрытый ключ и цепочку сертификатов на локальный диск.
Затем выполняется shell скрипт в котором вы можете скопировать файлы в нужную директорию и перезапустить http сервер
### Настройка
В файле settings.json укажите имкна сертификата, закрытого ключа и полной цепочки сертификатов а так же URL по которому они доступны
Например так:
        
        {
            "CertFile":"cert.cer",
            "KeyFile":"key.key",
            "FullChain": "fullchain.cer",
            "Url":"https://example.local/wafdirectory/certs"
        }

Впишите утилиту в crontab например так:

        0 19 * * sun cd /opt/getcert/ && ./GetCerts

### Сборка
go build ./GetCerts.go
либо выполняйте готовый бинарный linux64 файл GetCerts

### Отладка
При запуске можно видеть что происходит и какие сертификаты были получены
Закрытый ключ не выводится на экран

### Shell скрипт
Он зависит от вашего HTTP сервера но концепция такова:
Перезаписать текущие файлы (сертификаты и закрытый ключ) и перезапустить сервер
Например так:

            #!/bin/sh
            echo "stop httpd"
            service httpd stop
            echo "copy files"
            cp yourcert.cer /etc/httpd/conf/cert.crt
            cp yourkey.key /etc/httpd/conf/key.key
            echo "start httpd"
            service httpd start
