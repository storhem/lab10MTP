import pytest
from fastapi.testclient import TestClient
from main import app

client = TestClient(app)


def test_benchmark_ping(benchmark):
    result = benchmark(client.get, "/ping")
    assert result.status_code == 200


def test_benchmark_get_items(benchmark):
    result = benchmark(client.get, "/items")
    assert result.status_code == 200


def test_benchmark_get_item_by_id(benchmark):
    result = benchmark(client.get, "/items/1")
    assert result.status_code == 200


def test_benchmark_memory_endpoint(benchmark):
    result = benchmark(client.get, "/memory")
    assert result.status_code == 200
    data = result.json()
    assert "rss_mb" in data
    assert data["rss_mb"] > 0
