package grpool

import (

)
var pool Pool

type Job func()

type Worker struct {
	JobChannel chan Job
	Stop chan struct{}
}

type Pool struct {
	Workers chan *Worker
	JobQueue chan Job
}

func NewPool (NumWorkers int,  JobQueueLen int) *Pool {
	workers := make(chan *Worker, NumWorkers)
	jobqueue := make(chan Job, JobQueueLen)

	for i := 0; i <= cap(workers); i++{
		worker := NewWorker()
		worker.start()
	}
	pool = Pool{
		Workers: workers,
		JobQueue: jobqueue,
	}
	go dispatch()

	return &pool
}

func dispatch(){
	for {
		select {
		case job := <-pool.JobQueue:
			w := <-pool.Workers
			w.JobChannel<-job
		}

	}
}

func NewWorker() *Worker{
	return &Worker{
		JobChannel: make(chan Job),
	}

}

func (w *Worker) start()  {
	go func() {
		for {
			pool.Workers<- w
			select {
			case job := <-w.JobChannel:
				job()
			}
		}
	}()
}