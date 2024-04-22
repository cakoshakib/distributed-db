import socket
from time import sleep
import time

HOST = "127.0.0.1"
LEADER = 8080
NODE1 = 8070
N = 1000

def send_message(msg, port):
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((HOST, port))
        s.sendall(msg)
        s.recv(1024)

# send_message(b"get bob table1 key1;", NODE1)

# add user
init_time = time.time()
for i in range(2):
    msg = f"cu bob{i};"
    send_message(msg.encode("ascii"), LEADER)
post_time = time.time()
# print("Time taken for create user operation", post_time - init_time)

# create tables
init_time = time.time()
for i in range(2):
    msg = f"ct bob{i} table{i};"
    send_message(msg.encode("ascii"), LEADER)
post_time = time.time()
# print("Time taken for create user operation", post_time - init_time)

# add time
init_time = time.time()
for i in range(N):
    msg = f"add bob1 table1 key{i} value{i};"
    send_message(msg.encode("ascii"), LEADER)

post_time = time.time()
print("Time taken for add operation", post_time - init_time)

# get time
init_time = time.time()
for i in range(N):
    msg = f"get bob1 table1 key{i};"
    send_message(msg.encode("ascii"), LEADER)
post_time = time.time()
print("Time taken for get operation", post_time - init_time)




