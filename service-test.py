import unittest
import requests
import json
import time


class TestService(unittest.TestCase):
    base_url = "http://localhost:8080"

    def test_valid_request(self):
        # Правильный запрос
        params = {
            "X1": 10.0,
            "X2": 2.0,
            "X3": 5.0,
            "Y1": 15.0,
            "Y2": 3.0,
            "Y3": 5.0,
            "E": 2
        }
        response = requests.post(f"{self.base_url}/calculate", json=params)
        self.assertEqual(response.status_code, 200)
        data = response.json()
        self.assertIn("X", data)
        self.assertIn("Y", data)
        self.assertIn("IsEqual", data)

    def test_invalid_request(self):
        # Запрос с невалидными параметрами (деление на ноль)
        params = {
            "X1": 10.0,
            "X2": 0.0,
            "X3": 5.0,
            "Y1": 15.0,
            "Y2": 0.0,
            "Y3": 5.0,
            "E": 2
        }
        response = requests.post(f"{self.base_url}/calculate", json=params)
        self.assertEqual(response.status_code, 400)
        data = response.json()
        self.assertIn("error", data)

    def test_limit_exceeded(self):
        limit = 5
        interval = 5
        # Превышение лимита запросов
        time.sleep(interval)
        params = {
            "X1": 10.0,
            "X2": 2.0,
            "X3": 5.0,
            "Y1": 15.0,
            "Y2": 3.0,
            "Y3": 5.0,
            "E": 2
        }
        # Отправляем больше запросов, чем установленный лимит за интервал времени

        for i in range(limit + 1):
            response = requests.post(f"{self.base_url}/calculate", json=params)
            if i < limit:
                self.assertEqual(response.status_code, 200)
            else:
                self.assertEqual(response.status_code, 402)
                data = response.json()
                self.assertIn("error", data)

        # Ждем интервал времени, чтобы счетчик запросов сбросился
        time.sleep(interval)

        # После сброса счетчика запросы должны обрабатываться снова
        response = requests.post(f"{self.base_url}/calculate", json=params)
        self.assertEqual(response.status_code, 200)


if __name__ == "__main__":
    unittest.main()
