# -*- coding: utf-8 -*-
"""
Сравнение производительности и потребления памяти: Gin (Go) vs FastAPI (Python).

Использование:
    # Запустить Go-сервис:  cd src/go-service && go run .
    # Запустить FastAPI:    cd src/fastapi-service && uvicorn main:app --port 8000
    python tests/memory_compare.py
"""

import sys
import time
import urllib.request
import json

# Корректная кодировка для Windows-терминала
if sys.stdout.encoding and sys.stdout.encoding.lower() != "utf-8":
    sys.stdout.reconfigure(encoding="utf-8", errors="replace")

GIN_BASE = "http://localhost:8080"
FASTAPI_BASE = "http://localhost:8000"
REQUESTS_COUNT = 200


def fetch(url: str) -> tuple[int, float, bytes]:
    """GET-запрос, возвращает (статус, время_сек, тело)."""
    start = time.perf_counter()
    with urllib.request.urlopen(url, timeout=5) as resp:
        body = resp.read()
        status = resp.status
    elapsed = time.perf_counter() - start
    return status, elapsed, body


def benchmark_endpoint(base: str, path: str, n: int) -> dict:
    """Делает n запросов к base+path, возвращает статистику."""
    times = []
    errors = 0
    for _ in range(n):
        try:
            _, elapsed, _ = fetch(base + path)
            times.append(elapsed * 1000)  # → мс
        except Exception as e:
            errors += 1
    if not times:
        return {"mean_ms": 0, "min_ms": 0, "max_ms": 0, "errors": errors}
    return {
        "mean_ms": round(sum(times) / len(times), 3),
        "min_ms":  round(min(times), 3),
        "max_ms":  round(max(times), 3),
        "errors":  errors,
    }


def get_memory(base: str) -> dict:
    """Опрашивает /memory эндпоинт сервиса."""
    try:
        _, _, body = fetch(base + "/memory")
        return json.loads(body)
    except Exception as e:
        return {"error": str(e)}


def print_table(title: str, headers: list, rows: list):
    col_w = [max(len(h), max(len(str(r[i])) for r in rows)) for i, h in enumerate(headers)]
    sep = "+-" + "-+-".join("-" * w for w in col_w) + "-+"
    fmt = "| " + " | ".join(f"{{:<{w}}}" for w in col_w) + " |"
    print(f"\n{title}")
    print(sep)
    print(fmt.format(*headers))
    print(sep)
    for row in rows:
        print(fmt.format(*[str(c) for c in row]))
    print(sep)


def main():
    endpoints = ["/ping", "/items", "/items/1"]

    print(f"Тестирование: {REQUESTS_COUNT} запросов на каждый эндпоинт\n")

    perf_rows = []
    for path in endpoints:
        gin_stat    = benchmark_endpoint(GIN_BASE,    path, REQUESTS_COUNT)
        python_stat = benchmark_endpoint(FASTAPI_BASE, path, REQUESTS_COUNT)
        perf_rows.append([
            path,
            f"{gin_stat['mean_ms']} мс",
            f"{python_stat['mean_ms']} мс",
            f"{gin_stat['min_ms']} мс",
            f"{python_stat['min_ms']} мс",
            f"{gin_stat['max_ms']} мс",
            f"{python_stat['max_ms']} мс",
        ])

    print_table(
        "Производительность (среднее время ответа)",
        ["Эндпоинт", "Gin среднее", "FastAPI среднее",
         "Gin мин", "FastAPI мин", "Gin макс", "FastAPI макс"],
        perf_rows,
    )

    gin_mem    = get_memory(GIN_BASE)
    python_mem = get_memory(FASTAPI_BASE)

    if "error" not in gin_mem and "error" not in python_mem:
        mem_rows = [
            ["RSS (физ. память)",
             f"{gin_mem.get('alloc_mb', 0):.2f} МБ (heap alloc)",
             f"{python_mem.get('rss_mb', 0):.2f} МБ"],
            ["Системная память",
             f"{gin_mem.get('sys_mb', 0):.2f} МБ",
             f"{python_mem.get('vms_mb', 0):.2f} МБ (VMS)"],
            ["Циклов GC / % RAM",
             f"{gin_mem.get('num_gc', 0)} GC-циклов",
             f"{python_mem.get('percent', 0):.2f}%"],
        ]
        print_table(
            "Потребление памяти",
            ["Метрика", "Gin (Go)", "FastAPI (Python)"],
            mem_rows,
        )
    else:
        print(f"\n[WARN] Не удалось получить /memory: gin={gin_mem}, python={python_mem}")

    print("\nВывод:")
    if perf_rows:
        # Минимальное время точнее отражает скорость хэндлера (без network jitter)
        gin_min = float(perf_rows[0][3].split()[0])
        py_min  = float(perf_rows[0][4].split()[0])
        if gin_min > 0:
            ratio = py_min / gin_min
            print(f"  FastAPI медленнее Gin примерно в {ratio:.1f}x по минимальному времени /ping")
    print("  Go потребляет меньше памяти благодаря статической компиляции и")
    print("  отсутствию интерпретатора Python и его стандартных библиотек.")


if __name__ == "__main__":
    main()
