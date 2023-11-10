# object_model

Сервис для обмена данными с БД Postgres SQL
Обмен данными сделано разными методами:
1. Команды в NATS по чтению, изменению и др.
2. DB CRUD операции - прямой обмен с БД
   (в каждой модели таблицы есть методы Read(), Update(), Create(), Save(), Delete(), Restore())
3. GRPC - обмен с БД по протоколу GRPC
   (сервис клиент которому надо обмениваться с БД подключается к сервису серверу sync_exchange, последний обменивается с БД)
4. NRPC - обмен с БД по протоколу NRPC
   (сервис клиент которому надо обмениваться с БД подключается к сервису NATS, который передаёт команды серверу sync_exchange, последний обменивается с БД)

Перед началом выполнения CRUD операций надо указать транспорт по которому будет происходить обмен (CRUD, GRPC, NRPC)
с помощью одной из команд:
InitCrudTransport_DB()
InitCrudTransport_GRPC()
InitCrudTransport_NRPC()
из модуля
"gitlab.aescorp.ru/dsp_dev/claim/common/object_model/pkg/crud_starter"

Также для начала работы должны быть заполнены переменные окружения:

1) для DB CRUD:
DB_HOST="10.1.9.23"
DB_NAME="claim"
DB_SCHEME="public"
DB_PORT="5432"
DB_USER="dev"
DB_PASSWORD=

2) для GRPC:
SYNC_SERVICE_HOST=10.1.9.150
SYNC_SERVICE_PORT=30031

3) для NRPC:
BUS_LOCAL_HOST="10.1.9.150"
BUS_LOCAL_PORT=4222

Для NRPC (GRPC) желательно сначала подключиться туда и в конце отключиться
nrpc_client.Connect()
defer nrpc_client.CloseConnection()
иначе код всё равно туда подключится, и не отключится в конце работы микросервиса.

Образец кода можно найти в тестовых файлах, например:
https://gitlab.aescorp.ru/dsp_dev/claim/common/object_model/-/tree/main/pkg/nrpc/nrpc_client/nrpc_employees


