# sync_exchange
Пакет для синхронного обмена данными через шину (NATS).

## Принцип обмена
Сервис пишет в определённый топик шины, подписывается на другой топик (генерируется исходя из названия сервиса и идентификатора пакета) и ожидает сообщение. 
Слушающий сервис вычитывает топик, после обработки пакета отправляет ответ.

## Команды для формирования пакетов
### SyncPackageToJSON
View SyncPackage as JSON string

### SyncPackageFromJSON
Make SyncPackage from JSON string

### MakeSyncCommand
Create SyncPackage as command package

### MakeSyncResult
Create SyncPackage as result package

### MakeSyncError
Create SyncPackage as error package

## Команды для обмена
### InitSyncExchange
Функция инициализации подключения к шине

### DeInitSyncExchange
Функция деинициализации подключения к шине

### SendMessage
Отправка сообщения в шину без ожидания ответа

### WaitMessage
Ожидание сообщения из определённого топика из одной очереди. Аналог direct в rabbit. Сообщения будут доставляться по очереди разным подписчикам.
Имя очереди будет дефолтное sync_exchange.

### QueueSubscribe
Ожидание сообщения из определённого топика из указанной очереди. Аналог direct в rabbit. Сообщения будут доставляться по очереди разным подписчикам.

### Subscribe
Ожидание сообщения из определённого топика. Аналог fanout в rabbit. Сообщения будут доставляться всем подписчикам.

### SendRequest
Отправка запроса с ожиданием ответа

### SendResponse
Отправка ответа на запрос

## Пример использования
```go
package main

import (
	"fmt"
	"log"
	"gitlab.aescorp.ru/dsp_dev/test_area/test_claim/pkg/sync_exchange"
	"gitlab.aescorp.ru/dsp_dev/test_area/test_claim/pkg/sync_exchange/sync_types"
)

func main()  {
	err := sync_exchange.InitSyncExchange("localhost", "service_name")
	if err != nil {
		log.Fatal(err)
	}

	params := make(map[string]string)
	params["something"] = "sometime"
	pack := sync_types.MakeSyncCommand("command_1", params)
	resp, err := sync_exchange.SendRequest("service_new", pack, 10000)

	fmt.Println(resp.Body.Result["state"])

	err = sync_exchange.DeInitSyncExchange()
	if err != nil {
		log.Fatal(err)
	}
}
```

## Формат пакетов
### Команда MakeSyncCommand
```json
{
  "head": {
    "destVer": "1.0",
    "sender": "service_1",
    "netID": "810c8afd-6f88-4670-acf1-248fec76f2eb",
    "created": "2022-08-16 10:00:44.292"
  },
  "body": {
    "command": "new_command",
    "params": {
      "key1": "value1",
      "key2": "value2",
      "key3": "value3"
    },
    "error": {
      "place": "",
      "code": 0,
      "message": ""
    }
  }
}
```
### Команда MakeSyncResult
```json
{
  "head": {
    "destVer": "1.0",
    "sender": "service_1",
    "netID": "c10386f7-1601-49e1-b4ca-a3d856b72af9",
    "created": "2022-08-16 10:01:32.888"
  },
  "body": {
    "result": {
      "key1": "value1",
      "key2": "value2",
      "key3": "value3"
    },
    "error": {
      "place": "",
      "code": 0,
      "message": ""
    }
  }
}
```
### Команда MakeSyncError
```json
{
  "head": {
    "destVer": "1.0",
    "sender": "service_1",
    "netID": "0b2e1b8e-f187-4b9d-bea8-23f951177933",
    "created": "2022-08-16 09:58:47.206"
  },
  "body": {
    "error": {
      "place": "error place",
      "code": 123,
      "message": "error message"
    }
  }
}
```
