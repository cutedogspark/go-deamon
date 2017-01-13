package gworker

var Service = WorkerManage{}

type WorkerManage struct {
	Status         bool
	WorkerCnt      int
	WorkQueueCnt   int
	WorkerQueue    chan chan Job
	WorkQueue      chan Job
	workQueueAbort chan bool
	Workers        []*Worker
}

func (c *WorkerManage) AddItem(item *Worker) []*Worker {
	c.Workers = append(c.Workers, item)
	return c.Workers
}

func InitWorker(worker int) WorkerManage {

	if worker > 0 {
		Service.WorkerCnt = worker
	}
	Service.WorkerQueue = make(chan chan Job, Service.WorkerCnt)
	Service.WorkQueue = make(chan Job, Service.WorkQueueCnt)
	Service.workQueueAbort = make(chan bool)
	for i := 0; i < Service.WorkerCnt; i++ {
		worker := NewWorker(i+1, Service.WorkerQueue)
		Service.AddItem(&worker)
	}
	return Service
}

func (w *WorkerManage) Start() {
	if w.Status {
		return
	}
	go func() {
		for _, x := range w.Workers {
			x.Start()
		}
		w.Status = true
		for {
			select {
			case work := <-w.WorkQueue:
				go func() {
					worker := <-w.WorkerQueue
					worker <- work
				}()
			case <-w.workQueueAbort:
				return
			}
		}
	}()
}

func (w *WorkerManage) Stop() {
	go func() {
		for _, x := range w.Workers {
			x.Stop()
		}
		w.workQueueAbort <- true
		w.Status = false
	}()
}

func (w *WorkerManage) PutWorkQueue(work Job) bool {
	if w.Status {
		w.WorkQueue <- work
		return true
	}
	return false
}
