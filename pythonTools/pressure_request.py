"""
压测脚本：支持并发调用排行榜系统的四个接口：
- update: /update-score (POST)
- rank:   /get-rank    (GET)
- top:    /get-top     (GET)
- range:  /get-range   (GET)
# 更新分数压测
python pressure_request.py --action update --total 50000 --concurrency 100

# 查询排名压测
python pressure_request.py --action rank  --total 20000 --concurrency 50

# 查询前 N 名
python pressure_request.py --action top --n 20 --total 20000 --concurrency 50

# 查询范围
python pressure_request.py --action range --before 50 --after 50 --total 20000 --concurrency 100
"""

import requests
import time
import random
import argparse
from concurrent.futures import ThreadPoolExecutor
from requests.adapters import HTTPAdapter

# 默认配置
DEFAULT_TOTAL = 50000
DEFAULT_CONCURRENCY = 100
LOG_INTERVAL = 1000

BASE_URL = "http://127.0.0.1:8080"
ENDPOINTS = {
    "update": "/update-score",
    "rank":   "/get-rank",
    "top":    "/get-top",
    "range":  "/get-range"
}

# 全局统计
success_count = 0
fail_count = 0

# HTTP session with connection pooling
session = requests.Session()
adapter = HTTPAdapter(pool_connections=DEFAULT_CONCURRENCY,
                      pool_maxsize=DEFAULT_CONCURRENCY,
                      max_retries=1,
                      pool_block=True)
session.mount("http://", adapter)
session.mount("https://", adapter)

player_prefix = "test"
index = 0

def generate_random_player_id():
    global index
    index += 1
    return f"{player_prefix}{index}"

def send_update(player_id):
    global success_count, fail_count
    data = {
        "playerId": player_id,
        "score": random.randint(1, 1000000),
        "timestamp": int(time.time() * 1000)
    }
    try:
        r = session.post(BASE_URL + ENDPOINTS["update"], data=data, timeout=5)
        r.raise_for_status()
        success_count += 1
    except Exception as e:
        print(f"update exception: {e}")
        fail_count += 1

def send_rank(player_id):
    global success_count, fail_count
    params = {"playerId": player_id}
    try:
        r = session.get(BASE_URL + ENDPOINTS["rank"], params=params, timeout=5)
        r.raise_for_status()
        success_count += 1
    except Exception as e:
        print(f"rank exception: {e}")
        fail_count += 1

def send_top(n):
    global success_count, fail_count
    params = {"n": n}
    try:
        r = session.get(BASE_URL + ENDPOINTS["top"], params=params, timeout=5)
        r.raise_for_status()
        success_count += 1
    except Exception as e:
        print(f"top exception: {e}")
        fail_count += 1

def send_range(player_id, before, after):
    global success_count, fail_count
    params = {
        "playerId": player_id,
        "before": before,
        "after": after
    }
    try:
        r = session.get(BASE_URL + ENDPOINTS["range"], params=params, timeout=5)
        r.raise_for_status()
        success_count += 1
    except Exception as e:
        print(f"range exception: {e}")
        fail_count += 1

def main():
    global success_count, fail_count

    parser = argparse.ArgumentParser(description="Leaderboard API压力测试")
    parser.add_argument("--action", required=True, choices=["update","rank","top","range"],
                        help="测试接口类型")
    parser.add_argument("--playerId", type=str, help="玩家ID (rank/range/update)")
    parser.add_argument("--n", type=int, default=10, help="获取 top N (top)")
    parser.add_argument("--before", type=int, default=5, help="range before")
    parser.add_argument("--after", type=int, default=5, help="range after")
    parser.add_argument("--total", type=int, default=DEFAULT_TOTAL, help="总请求数")
    parser.add_argument("--concurrency", type=int, default=DEFAULT_CONCURRENCY, help="并发线程数")
    args = parser.parse_args()

    total = args.total
    concurrency = args.concurrency

    print(f"Action: {args.action}, total: {total}, concurrency: {concurrency}")
    start = time.time()

    with ThreadPoolExecutor(max_workers=concurrency) as pool:
        for i in range(total):
            # generate needed parameters
            if args.action == "update":
                pid = args.playerId or generate_random_player_id()
                pool.submit(send_update, pid)
            elif args.action == "rank":
                pid = args.playerId or generate_random_player_id()
                pool.submit(send_rank, pid)
            elif args.action == "top":
                pool.submit(send_top, args.n)
            elif args.action == "range":
                pid = args.playerId or generate_random_player_id()
                pool.submit(send_range, pid, args.before, args.after)
            # log progress
            if (i+1) % LOG_INTERVAL == 0:
                print(f"dispatched: {i+1}/{total}")

    duration = time.time() - start
    print("\nTest End")
    print(f"Success: {success_count}, Fail: {fail_count}")
    print(f"Duration: {duration:.2f}s, QPS: {total/duration:.2f}")

if __name__ == "__main__":
    main()

