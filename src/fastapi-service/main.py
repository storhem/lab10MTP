import os

import psutil
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

app = FastAPI(
    title="Items API",
    description=(
        "REST API для управления каталогом товаров. "
        "Часть лабораторной работы №10 — сравнение FastAPI (Python) и Gin (Go)."
    ),
    version="1.0.0",
    contact={
        "name": "Евланичев Максим Юрьевич",
        "email": "storhetmax@mail.ru",
    },
)


class Item(BaseModel):
    id: int
    name: str
    price: float


class MemoryStats(BaseModel):
    rss_mb: float
    vms_mb: float
    percent: float


class MessageResponse(BaseModel):
    message: str


class ErrorResponse(BaseModel):
    detail: str


ITEMS: list[Item] = [
    Item(id=1, name="Apple", price=1.5),
    Item(id=2, name="Banana", price=0.75),
]


@app.get(
    "/memory",
    response_model=MemoryStats,
    tags=["metrics"],
    summary="Потребление памяти",
    description="Возвращает метрики памяти текущего процесса через psutil.",
)
def memory():
    proc = psutil.Process(os.getpid())
    mem = proc.memory_info()
    return MemoryStats(
        rss_mb=mem.rss / 1024 / 1024,
        vms_mb=mem.vms / 1024 / 1024,
        percent=proc.memory_percent(),
    )


@app.get(
    "/ping",
    response_model=MessageResponse,
    tags=["health"],
    summary="Проверка работоспособности",
    description="Возвращает `pong` — используется для health-check.",
)
def ping():
    return {"message": "pong"}


@app.get(
    "/items",
    response_model=list[Item],
    tags=["items"],
    summary="Список всех товаров",
    description="Возвращает полный список товаров в каталоге.",
)
def get_items():
    return ITEMS


@app.get(
    "/items/{item_id}",
    response_model=Item,
    tags=["items"],
    summary="Товар по ID",
    description="Возвращает товар по его числовому идентификатору.",
    responses={404: {"model": ErrorResponse, "description": "Товар не найден"}},
)
def get_item(item_id: int):
    for item in ITEMS:
        if item.id == item_id:
            return item
    raise HTTPException(status_code=404, detail="item not found")
