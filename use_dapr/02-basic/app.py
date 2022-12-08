import os

from fastapi import FastAPI
import uvicorn

PORT = int(os.environ.get("PORT", 5000))


app = FastAPI()


@app.post("/add")
def add():
    return "Added!"


@app.get("/index")
def index():
    return f"index app: {PORT}"


if __name__ == "__main__":
    uvicorn.run(app, port=PORT)
