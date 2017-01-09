
# go-deamon

## gworker: A Golang Worler & Job Package.


## Environment

- [x] go 1.7.4
- [x] glide 0.12.3
- [x] echo 3.0.3

## Installation


## Sample Code in ``example\`` folder
```sh
./gworker_example
 ⇛ http server started on 127.0.0.1:8000
Starting worker(4) ......Success
```
Run test command
```sh
for i in {1..4}; do curl localhost:8000/work -d name=$USER -d delay=$(expr $i % 11)s; done
```


## Usage
```go
var QWorker gworker.WorkerManage
...
QWorker = gworker.InitWorker(*NWorkers)
QWorker.Start()
...
```

```go
work := gworker.Job{Name: name, Delay: delay}
QWorker.PutWorkQueue(work)
```


Hava fun!