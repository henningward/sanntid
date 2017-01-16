from threading import Thread
from threading import Lock

i = 0
lock = Lock()

def firstThreadFunc():
    global i
    for j in range(0, 1000000):
        lock.acquire()
        i+=1
        lock.release()

def secondThreadFunc():
    global i
    for j in range(0, 1000000):
        lock.acquire()
        i-=1
        lock.release()


# Potentially useful thing:
#   In Python you "import" a global variable, instead of "export"ing it when you declare it
#   (This is probably an effort to make you feel bad about typing the word "global")




def main():
    global i
    FirstThread= Thread(target = firstThreadFunc, args = (),)
    FirstThread.start()

    SecondThread= Thread(target = secondThreadFunc, args = (),)
    SecondThread.start()
    
    FirstThread.join()
    SecondThread.join()
    print(i)
	
main()

