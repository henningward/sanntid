#include <pthread.h>
#include <stdio.h>
#define num_threads 2

extern int i = 0;


void* thread_function(void* j){
	if (*((int*)j) == 0){
		for (int k = 0; k < 1000000; k++){
			i++;
		}
	} else if (*((int*)j) == 1){
		for (int k = 0; k < 1000000; k++){
			i--;
		}
	}
	return NULL;
}


int main(){

	int j = 0;

	pthread_t thread[num_threads];
	for (j = 0; j < num_threads; j++){
		pthread_create(&thread[j], NULL, thread_function, &j);

	}

    for (j = 0; j < num_threads; j++){
    	pthread_join(thread[j], NULL);
    }


	// pthread_t thread;
	// pthread_create(&thread, NULL, thread_function, &j);
	// pthread_join(thread, NULL);
    printf("%d\n", i);
    return 0;
    
}