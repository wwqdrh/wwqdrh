from fastapi import FastAPI
import uvicorn

app = FastAPI()


@app.post("/add")
def add():
    return "Added!"

@app.get("/index")
def index():
    return "app1"


if __name__ == "__main__":
    uvicorn.run(app, host='0.0.0.0', port=8080)
