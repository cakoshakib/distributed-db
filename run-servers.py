import sys
import subprocess
import os
from multiprocessing.pool import ThreadPool
from multiprocessing import Lock
import time

TCP_PORT = 8000
LEADER_PORT = TCP_PORT
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
    print()
    print(f"Running server{i}")
    print(f"    - NAME: node{i}")
    print(f"    - DATA: ./data/server{i}/") 
    print(f"    - TCP PORT: {TCP_PORT + i}")
    print(f"    - RAFT PORT: {RAFT_PORT + i}")
    print(f"    - LOG: ./logs/server{i}.log", flush=True)
    print()

def run(i):
    global procs
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
    if len(procs) != 0: 
        # TODO: Find a joinAddr workaround
        # Theres no good way to find the LEADER_PORT, as Raft randomly assigns a leader.
        cmd += ["-joinAddr", f"localhost:{LEADER_PORT}"]
    with open(f"./logs/server{i}.log", "w+") as logf:
        ret = subprocess.Popen(cmd, stdout=logf, stderr=logf)
    return ret


procs = {}
max_n = num_servers - 1
for i in range(num_servers):
    node_id = f"node{i}"
    procs[node_id] = run(i)
    time.sleep(2)

"""
help
list
kill i
add
"""
def help():
    print()
    print("Available commands:")
    print("  - list : lists all nodes running")
    print("  - kill [i] : kills node i from the cluster")
    print("  - add : adds node to cluster and prints info")
    print("  - quit : kill entire cluster")
    print()

def list():
    print()
    print("List of running nodes:")
    for key in procs:
        print(f"  - {key}")
    print()

def kill(i):
    print()
    print(f"Attempting to kill node {i}...")
    killed = False
    if i in procs:
        proc = procs[i]
        proc.kill()
        print(f"Killed node {i}")
        killed = True
        del procs[i]
    if not killed:
        print(f"Node {i} does not exist")
    print()

def add():
    print()
    global max_n
    max_n += 1
    node_id = f"node{max_n}"
    print(f"Adding node {max_n}")
    procs[node_id] = run(max_n)

def quit():
    for _, proc in procs:
        proc.kill()

while True:
    print(">", end="", flush=True)
    command = input().lower()
    if command == "help":
        help()
    elif command == "list":
        list()
    elif command.startswith("kill"):
        parts = command.split()
        if len(parts) < 2:
            print("Provide index to kill")
        else:
            kill(parts[1])

    elif command == "add":
        add()
    elif command == "quit":
        quit()
        break
