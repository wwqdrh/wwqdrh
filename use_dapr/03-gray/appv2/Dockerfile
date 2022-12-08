FROM python:3.10

WORKDIR app

RUN pip install fastapi uvicorn -i https://pypi.tuna.tsinghua.edu.cn/simple

COPY ./app.py app.py

ENTRYPOINT ["python", "/app/app.py"]
