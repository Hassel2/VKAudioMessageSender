# VKAudioMessageSender

Программа которая отправляет голосовые сообщения из локального аудиофайла.

Инструкция по установке:
1) Для начала нужно установить Golang последней версии.
2) После этого необходимо клонировать репозиторий
3) Переходим в дерикторию репозитория и выполняем команду go install, которая автоматически соберет бинарный файл и положит его в ~/go/bin
4) После этого нужно добавить ~/go/bin в PATH
5) На этом все.

Инструкция по использованию:

В терминале пишем:
bassboost -p /path/to/audiofile -id 123456 (id получаетля)
