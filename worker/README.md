# Asynchronous tasks

One way to implement asynchronous tasks is go routine (background thread in process).
It's easy to implement but executed in the same process and if there're some problem, it can lose the unprocess tasks.

Other way to do that is [Redis](https://redis.io/), message broker and background worker.

- Tasks saved in both memory and persistent storage
- Highly available: Redis sentinel and Redis cluster
- No task lost

## Asynq package

Tasks typically managed by queue and worker, [Asynq](https://github.com/hibiken/asynq) is one of such management library.

There're a distributor and a processor for the task management

- Distributor: Enqueue tasks
- Processor: Process each task

see [distributor.go](./distributor.go) and [processor.go](./processor.go).
The processor will be stand alone server to do them.
