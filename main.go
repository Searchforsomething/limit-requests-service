package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
)

// Config содержит конфигурацию сервиса
type Config struct {
	Limit    int           // Лимит на количество запросов
	Interval time.Duration // Интервал времени для лимита
}

// RequestCounter отслеживает количество запросов
type RequestCounter struct {
	count     int64
	lastReset int64
	Limit     int64
	Interval  int64 // Интервал времени для лимита в наносекундах
}

// NewRequestCounter создает новый счетчик запросов с заданным лимитом и интервалом
func NewRequestCounter(limit int, interval time.Duration) *RequestCounter {
	return &RequestCounter{
		count:     0,
		lastReset: time.Now().UnixNano(),
		Limit:     int64(limit),
		Interval:  interval.Nanoseconds(),
	}
}

// Inc увеличивает счетчик запросов. Возвращает true, если лимит не превышен, иначе false.
func (rc *RequestCounter) Inc() bool {
	now := time.Now().UnixNano()

	if now-atomic.LoadInt64(&rc.lastReset) > rc.Interval {
		atomic.StoreInt64(&rc.count, 0)
		atomic.StoreInt64(&rc.lastReset, now)
	}

	if atomic.LoadInt64(&rc.count) < rc.Limit {
		atomic.AddInt64(&rc.count, 1)
		return true
	}
	return false
}

// Parameters содержит входные параметры запроса
type Parameters struct {
	X1 float64 `json:"X1"`
	X2 float64 `json:"X2"`
	X3 float64 `json:"X3"`
	Y1 float64 `json:"Y1"`
	Y2 float64 `json:"Y2"`
	Y3 float64 `json:"Y3"`
	E  int     `json:"E"` // Точность (количество знаков после запятой)
}

// Response содержит ответ сервиса
type Response struct {
	X       float64 `json:"X"`
	Y       float64 `json:"Y"`
	IsEqual string  `json:"IsEqual"`
}

var counter *RequestCounter

func main() {
	// Чтение параметров командной строки
	limit := flag.Int("limit", 5, "Лимит запросов")
	interval := flag.Duration("interval", 5*time.Second, "Интервал времени для лимита")
	flag.Parse()

	// Создание счетчика запросов
	counter = NewRequestCounter(*limit, *interval)

	// Обработчик POST запросов
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		// Проверяем лимит запросов
		if !counter.Inc() {
			ctx.SetStatusCode(402)
			fmt.Fprintf(ctx, `{"error": "Превышен лимит запросов. Пожалуйста, попробуйте позже."}`)
			return
		}

		// Чтение тела запроса
		var params Parameters
		err := json.Unmarshal(ctx.PostBody(), &params)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			fmt.Fprintf(ctx, `{"error": "Ошибка чтения параметров запроса"}`)
			return
		}

		// Проверка корректности входных данных
		if err := validateParameters(params); err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			fmt.Fprintf(ctx, `{"error": "%s"}`, err.Error())
			return
		}

		// Выполнение вычислений
		X, err := calculate(params.X1, params.X2, params.X3, params.E)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			fmt.Fprintf(ctx, `{"error": "Ошибка вычислений X: %s"}`, err.Error())
			return
		}

		Y, err := calculate(params.Y1, params.Y2, params.Y3, params.E)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusInternalServerError)
			fmt.Fprintf(ctx, `{"error": "Ошибка вычислений Y: %s"}`, err.Error())
			return
		}

		isEqual := "F"
		if X == Y {
			isEqual = "T"
		}

		// Формирование ответа
		resp := Response{
			X:       X,
			Y:       Y,
			IsEqual: isEqual,
		}

		// Возвращаем ответ в формате JSON
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusOK)
		responseBody, _ := json.Marshal(resp)
		ctx.SetBody(responseBody)
	}

	// Запуск HTTP сервера
	log.Fatal(fasthttp.ListenAndServe(":8080", requestHandler))
}

// validateParameters проверяет корректность входных параметров
func validateParameters(params Parameters) error {
	if params.E < 0 {
		return fmt.Errorf("значение E должно быть неотрицательным")
	}
	if params.X2 == 0 {
		return fmt.Errorf("деление на ноль: X2 равен нулю")
	}
	if params.Y2 == 0 {
		return fmt.Errorf("деление на ноль: Y2 равен нулю")
	}
	return nil
}

// calculate вычисляет значение на основе входных параметров
func calculate(a, b, c float64, precision int) (float64, error) {
	result := a / b * c
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, fmt.Errorf("некорректный результат вычислений")
	}
	return round(result, precision), nil
}

// round округляет число f до prec знаков после запятой
func round(f float64, prec int) float64 {
	pow := math.Pow10(prec)
	return math.Round(f*pow) / pow
}
