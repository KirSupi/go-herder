import sys
import time

print(sys.argv[1], end="")
while True:
    time.sleep(3)
    print(sys.argv[1], "...")
