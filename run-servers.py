import sys
import subprocess
import os

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

# create log folder if it does not exist
if not os.path.exists("./log"):
    print("Creating log directory")
    os.mkdir("./log")

# run each server
for i in range(num_servers):
    subprocess.run(["./distributed-db"])

