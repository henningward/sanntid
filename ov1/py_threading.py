from threading import Thread


i = 0
def firstThreadFunc():
	global i
	for j in range(0, 1000000):
		i+=1


def secondThreadFunc():
	global i
	for j in range(0, 1000000):
	   	i-=1

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

