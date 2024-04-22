import sys
import subprocess
import os
from multiprocessing.pool import ThreadPool
from multiprocessing import Lock
import time

TCP_PORT = 8000
RAFT_PORT = 12000

# read number of servers
if len(sys.argv) <= 1:
    print("USAGE: python run-servers.py [num_servers]")
    exit(0)
num_servers = int(sys.argv[1])

# build go program
print("Building go program")
subprocess.run(["go", "build"], capture_output=True)

# create relevant directories
def create_directory(path):
    if not os.path.exists(path):
        print(f"Creating {path} directory")
        os.mkdir(path)

create_directory("./logs")
create_directory("./data")
create_directory("./boltdbstore")

# run each server
lock = Lock()

def print_server_info(i):
    print(f"Running server{i}")
    print(f"    - NAME: node{i}")
    print(f"    - DATA: ./data/server{i}/") 
    print(f"    - TCP PORT: {TCP_PORT + i}")
    print(f"    - RAFT PORT: {RAFT_PORT + i}")
    print(f"    - LOG: ./logs/server{i}.log")

def run(i):
    lock.acquire()
    print_server_info(i)
    create_directory(f"./data/server{i}")
    cmd = ["./distributed-db", 
        "-tcpPort", str(TCP_PORT + i), 
        "-id", f"node{i}",
        "-raftAddr", f"localhost:{RAFT_PORT + i}",
        "-dataDir", f"data/server{i}/"
    ]
    lock.release()
    # join if not the first node
    if i != 0: 
        cmd += ["-joinAddr", f"localhost:{TCP_PORT}"]
    with open(f"./logs/server{i}.log", "w+") as logf:
        subprocess.run(cmd, stdout=logf, stderr=logf)


with ThreadPool(num_servers) as pool:
    pool.map_async(run, range(0,1))
    time.sleep(3)
    pool.map(run, range(1,num_servers))


