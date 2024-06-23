# Вычисления и лимитирование числа запросов

## Ограничения сервиса

1. **Безопасность входных данных:** 
сейчас не производится достаточной валидации входных данных.
2. **Масштабируемость:** текущая реализация работает в пределах одного процесса.
Для высоконагруженных систем требуется решение на уровне кластера.
3. **Мониторинг и логирование:** в реализации этого сервиса отсутствуют
механизмы логирования, которые необходимы для оперативного анализа проблем

## Структура проекта:

- `main.go` - основной код сервиса
- `service-test.py` - функциональные тесты сервиса

## Инструкция по запуску

Для запуска сервиса необходимо выполнить в терминале:
```bash
go run main.go -limit=<лимит запросов> -interval=<интервал времени в секундах>s
```
При отсутствии явного указания значений флагов сервис принимает значения по 
умолчанию: 
- лимит запросов = 5
- интервал времени = 5 секунд
  
В теле запроса к сервису передаются данные в формате json со значениями X1, X2, X3, Y1, Y2, Y3, E.  
Пример запроса с использованием curl:
```bash
curl -X POST -d '{"X1": 10.4, "X2": 2.1, "X3": 2.5, "Y1": 8.2, "Y2": 1.4, "Y3": 1.5, "E": 2}' http://localhost:8080/calculate
```