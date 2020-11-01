# priceServer
> *Сервис, позволяющий следить за изменением цены любого объявления на Авито.* 

### Запуск
```bash
$ docker-compose up
```

### Подписка на рассылку
Принимает почту по ключу `Mail`, ссылку по ключу `Link`.
Пример почты: mail@yandex.ru
```bash
$ curl --header "Content-Type: application/json" --request POST --data '{"Mail":"mail@domain","Link":"https://www.avito.ru/link"}' http://localhost:9090/subscribe
```

### Почта, используемся в рассылке
Отправление писем происходит с почты, которую надо указать в `config.txt`\
В поле `mail` надо указать почту. В поле `password` - пароль.
