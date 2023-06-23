# write an http post in python
import requests
import sys

# function that generates a random string and returns that string
def random_string():
    import random
    import string
    return ''.join(random.choice(string.ascii_lowercase) for i in range(10))

# curl -X POST http://localhost:10000/update -d '{"key":"test","data":"{\"something\":\"yes\"}"}'
if __name__ == "__main__":
    for x in range(100):
        keyValue = random_string()
        data = '{"key":"'+keyValue+'","data":"test"}'
        # print(data)
        r = requests.post("http://192.168.1.19:31048/update", data=data)
        # print(r.text)
        for y in range(10):
            r = requests.post("http://192.168.1.19:31048/get/"+keyValue, data=data)
            # print(r.text)
