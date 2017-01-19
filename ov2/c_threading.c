#include <pthread.h>
#include <stdio.h>
#define num_threads 2

//gcc -std=gnu99 -Wall -g -o c_threading c_threading.c -lpthread;
//./c_threading

int  i = 0;
pthread_mutex_t mutex;

void* thread_function(void* j){
	if (*((int*)j) == 0){
		for (int k = 0; k < 1000000; k++){
			pthread_mutex_lock(&mutex);

			i++;
    		pthread_mutex_unlock(&mutex);
	        
		}
	}
    else if (*((int*)j) == 1){
		for (int k = 0; k < 1000000; k++){
			pthread_mutex_lock(&mutex);
			i--;
    		pthread_mutex_unlock(&mutex);
	        
		}
	}
    return NULL;
}


int main(){
    pthread_mutex_init(&mutex, NULL);
    int j = 0;

	pthread_t thread[num_threads];
	for (j = 0; j < num_threads; j++){
        //printf(" j = %d", j);
		pthread_create(&thread[j], NULL, thread_function, &j);

	}

    for (int l = 0; l < num_threads; l++){
    	pthread_join(thread[l], NULL);
    }


	// pthread_t thread;
	// pthread_create(&thread, NULL, thread_function, &j);
	// pthread_join(thread, NULL);
    printf("%d\n", i);
    
    return 0;
    
}
