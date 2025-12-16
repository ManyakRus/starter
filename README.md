Набор компонент для языка Golang
Автор: Александр Никитин

Набор компонент для языка golang сделан для облегчения работы программиста,
чтобы любой компонент можно было подключить одной строкой кода.

Компоненты для запуска любых микросерверов (ядро):
1. log - компонент для логирования информации в консоль (логгер logrus)
2. contextmain - контекст общий на всё приложение
3. stopapp - ожидание завершения работы приложения Gracefull shutdown (CTRL+C), WaitGroup
4. config - загрузка параметров из файла .env или из переменных окружения

Подключение к внешним сервисам:
1. camunda_connect - подключение с сервису camunda
2. chatgpt_connect - подключение к сервису ChatGPT OpenAI, искуственный интеллект
3. fiber_connect - подключение веб сервера с компонентой fiber
4. kafka_connect - подключение к брокеру сообщений kafka
5. liveness - создание примитивного веб сервера для проверки работает или нет микросервис
6. mssql_connect - подключение к серверу Microsoft SQL server с драйвером sqlx
7. mssql_gorm - подключение к серверу Microsoft SQL server с драйвером gorm
8. nats_connect - подключение к брокеру сообщений NATS
9. postgres_connect - подключение с серверу баз данных Postgres, с драйвером sqlx
10. postgres_gorm - подключение с серверу баз данных Postgres, с драйвером gorm
11. postgres_pgx - подключение с серверу баз данных Postgres, с драйвером pgx
12. whatsapp_connect - подключение к сервисам мессенджера whatsapp

Каждое подключение к внешним сервисам использует общий logger, contextmain, WaitGroup,
config, и свою структуру Settings с параметрами

А также дополнительные библиотеки:
1. micro - набор небольших полезных функций
2. ping - функция для проверки работы порта на нужном хосте

Пример минимального ядра:
```
func main() {
	config.LoadEnv()
	stopapp.StartWaitStop()

	//ваш код

	stopapp.Wait_GracefulShutdown()
}
```


Пример с разными подключениями:
```
func main() {
	config.LoadEnv()

	contextmain.GetContext()

	stopapp.StartWaitStop()

	mssql_gorm.StartDB()

	postgres_gorm.StartDB()

	nats.StartNats()

	camunda.StartCamunda()

	liveness.Start()

	stopapp.Wait_GracefulShutdown()

	log.Info("App stopped")
}
```
