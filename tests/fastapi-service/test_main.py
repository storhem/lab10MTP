import pytest
from fastapi.testclient import TestClient
from main import app

client = TestClient(app)


def test_ping():
    response = client.get("/ping")
    assert response.status_code == 200
    assert response.json() == {"message": "pong"}


def test_get_items():
    response = client.get("/items")
    assert response.status_code == 200
    data = response.json()
    assert isinstance(data, list)
    assert len(data) > 0
    assert data[0]["name"] == "Apple"


def test_get_item_by_id_found():
    response = client.get("/items/1")
    assert response.status_code == 200
    assert response.json()["id"] == 1


def test_get_item_by_id_not_found():
    response = client.get("/items/999")
    assert response.status_code == 404
    assert "detail" in response.json()


def test_get_item_by_id_invalid():
    response = client.get("/items/abc")
    assert response.status_code == 422


def test_memory_endpoint():
    response = client.get("/memory")
    assert response.status_code == 200
    data = response.json()
    assert "rss_mb" in data
    assert "vms_mb" in data
    assert "percent" in data
    assert data["rss_mb"] > 0


def test_swagger_ui_available():
    response = client.get("/docs")
    assert response.status_code == 200
    assert "swagger" in response.text.lower()


def test_openapi_schema_available():
    response = client.get("/openapi.json")
    assert response.status_code == 200
    schema = response.json()
    assert schema["info"]["title"] == "Items API"
    assert schema["info"]["version"] == "1.0.0"
    assert "/ping" in schema["paths"]
    assert "/items" in schema["paths"]
    assert "/items/{item_id}" in schema["paths"]
