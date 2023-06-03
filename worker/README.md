# Asynchronous tasks

One way to implement asynchronous tasks is go routine (background thread in process).
It's easy to implement but executed in the same process and if there're some problem, it can lose the unprocess tasks.

Other way to do that is [Redis](https://redis.io/), message broker and background worker.

- Tasks saved in both memory and persistent storage
- Highly available: Redis sentinel and Redis cluster
- No task lost

## Asynq package

Tasks typically managed by queue and worker, [Asynq](https://github.com/hibiken/asynq) is one of such management library.

Note: asynq is still under heavy development and there might be breaking changes.

There're a distributor and a processor for the task management

- Distributor: Enqueue tasks when called
    ```go
    task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
    ```
- Processor: Process each task on the other process (or server)
    ```go
    mux := asynq.NewServeMux()
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)
	return processor.server.Start(mux)
    ```

see [distributor.go](./distributor.go) and [processor.go](./processor.go).
The processor will be stand alone server to do them.

### Redis Config and Run Server

Redis server can run with [redis-docker](https://hub.docker.com/_/redis).
After that, the task processor redis client is up, [main.go](../main.go).

To distribute task, we use the [distributor](./distributor.go) like [rpc_create_user.go](../gapi/rpc_create_user.go).
```go
taskPayload := &worker.PayloadSendVerifyEmail{
    Username: user.Username,
}
opts := []asynq.Option{
    asynq.MaxRetry(10),                // up to 10 retry
    asynq.ProcessIn(10 * time.Second), // 10 seconds delay
    asynq.Queue(worker.QueueCritical), // add critical instead of default
}
err = server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
if err != nil {
    return nil, status.Errorf(codes.Internal, "failed to distribute task to send verify email %s", err)
}
```

One advantage to use redis is an easy configuration.
```go
server := asynq.NewServer(
    redisOpt,
    asynq.Config{
        Queues: map[string]int{
            QueueCritical: 10,
            QueueDefault:  5,
        },
    },
)
```

### Redis Error Handler and logger

For custom logging of redis,

```go
type Logging interface {
	Printf(ctx context.Context, format string, v ...interface{})
}
```

is need to implement.
Similarly, asynq also has same strategy for custom logging
and [asynq one](https://github.com/hibiken/asynq/blob/master/server.go#L268) is

```go
// Logger supports logging at various log levels.
type Logger interface {
	// Debug logs a message at Debug level.
	Debug(args ...interface{})

	// Info logs a message at Info level.
	Info(args ...interface{})

	// Warn logs a message at Warning level.
	Warn(args ...interface{})

	// Error logs a message at Error level.
	Error(args ...interface{})

	// Fatal logs a message at Fatal level
	// and process will exit with status set to 1.
	Fatal(args ...interface{})
}
```

see [logger.go](./logger.go).

To handle errors in redis worker, we must set error handler in the config. Otherwise, nothing will happen by default.

```go
server := asynq.NewServer(
    redisOpt,
    asynq.Config{
        Queues: map[string]int{
            QueueCritical: 10,
            QueueDefault:  5,
        },
        ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
            log.Error().Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("process task failed")
        }),
        Logger: logger,
    },
)
```
