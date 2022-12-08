import time
import os

LOG = os.getenv("LOG", "/data/log/test.log")

if __name__ == "__main__":
    times = 0
    while True:
        times+=1
        print(f"====第{times}尝试读取日志")
        with open(LOG, mode="r", encoding="utf8") as f:
            print(f.read())
        time.sleep(5)
