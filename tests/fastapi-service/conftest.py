import sys
import os

# Добавляем src/fastapi-service в путь один раз для всех тестов пакета.
sys.path.insert(0, os.path.join(os.path.dirname(__file__), "../../src/fastapi-service"))
