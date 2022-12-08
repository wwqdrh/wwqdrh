import uvicorn
import fastapi
import pydantic


class ModelUser(pydantic.BaseModel):
    id: int
    name: str
    password: str


_id = 0


def GetID() -> int:
    global _id
    _id += 1
    return _id


MockUserDB: dict[str, ModelUser] = dict()

app = fastapi.FastAPI()


class ReqRegister(pydantic.BaseModel):
    name: str
    password: str


@app.post("/user/register")
async def register(data: ReqRegister):
    MockUserDB[data.name] = ModelUser(
        id=GetID(), name=data.name, password=data.password
    )
    return {"code": 0, "msg": "注册成功"}


class ReqLogin(pydantic.BaseModel):
    name: str
    password: str


@app.post("/user/login")
async def login(data: ReqLogin):
    if data.name not in MockUserDB:
        return {"code": 1, "msg": "用户名不存在"}
    else:
        return {"code": 0, "msg": "登录成功", "access_token": "123456"}


@app.post("/user/info")
async def info(token: str = fastapi.Header(alias="x-token")):
    if token != "123456":
        return {"code": 1, "msg": "token错误"}
    elif len(MockUserDB) == 0:
        return {"code": 1, "msg": "无数据"}
    else:
        # a random user
        
        return {
            "code": 0,
            "msg": "查询成功",
            "data": MockUserDB[next(iter(MockUserDB.keys()))],
        }


if __name__ == "__main__":
    uvicorn.run("how_api_test.api:app", port=8081, host="0.0.0.0")
