#!/usr/bin/env bash
# Нагрузочное тестирование: Gin vs FastAPI
# Инструмент: hey (https://github.com/rakyll/hey)
# Использование: bash benchmark.sh

HEY=$(go env GOPATH)/bin/hey

GIN_URL="http://localhost:8080"
FASTAPI_URL="http://localhost:8000"

echo "Сборка Go-сервиса..."
(cd src/go-service && go build -o bin/go-service.exe .)

echo "Запуск Go-сервиса..."
GIN_MODE=release src/go-service/bin/go-service.exe &
GIN_PID=$!

echo "Запуск FastAPI-сервиса..."
# Запуск из директории сервиса — Python не может импортировать пакеты с дефисом
(cd src/fastapi-service && python -m uvicorn main:app --host 0.0.0.0 --port 8000) &
FASTAPI_PID=$!

sleep 3

run_bench() {
    local label=$1
    local url=$2
    local n=$3
    local c=$4
    echo ""
    echo "=============================="
    echo "  $label | n=$n c=$c"
    echo "=============================="
    $HEY -n "$n" -c "$c" "$url"
}

# Сценарий 1: лёгкая нагрузка
run_bench "GIN    /ping" "$GIN_URL/ping"     1000 50
run_bench "FASTAPI /ping" "$FASTAPI_URL/ping" 1000 50

# Сценарий 2: JSON-ответ
run_bench "GIN    /items" "$GIN_URL/items"     1000 50
run_bench "FASTAPI /items" "$FASTAPI_URL/items" 1000 50

# Сценарий 3: высокая нагрузка
run_bench "GIN    /ping" "$GIN_URL/ping"     5000 100
run_bench "FASTAPI /ping" "$FASTAPI_URL/ping" 5000 100

kill $GIN_PID $FASTAPI_PID 2>/dev/null
echo ""
echo "Готово."
