import time
import os

LOG = os.getenv("LOG", "/data/log/test.log")
APP = os.getenv("APP", "app")

if __name__ == "__main__":
    times = 0
    while True:
        times+=1
        print(f"====第{times}尝试追加日志")
        with open(LOG, mode="a+", encoding="utf8") as f:
            f.write(f"[{APP}]: a {times} line log\n")
        time.sleep(5)
