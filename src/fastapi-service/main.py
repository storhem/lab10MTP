from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

app = FastAPI(title="FastAPI Service", version="1.0.0")


class Item(BaseModel):
    id: int
    name: str
    price: float


ITEMS: list[Item] = [
    Item(id=1, name="Apple", price=1.5),
    Item(id=2, name="Banana", price=0.75),
]


@app.get("/ping")
def ping():
    return {"message": "pong"}


@app.get("/items", response_model=list[Item])
def get_items():
    return ITEMS


@app.get("/items/{item_id}", response_model=Item)
def get_item(item_id: int):
    for item in ITEMS:
        if item.id == item_id:
            return item
    raise HTTPException(status_code=404, detail="item not found")
