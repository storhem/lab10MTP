import sys
import os

sys.path.insert(0, os.path.join(os.path.dirname(__file__), "../../src/fastapi-service"))

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
