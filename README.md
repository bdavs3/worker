# Worker

**Warning**: This README describes proposed usage of the service, which is still in development.

`worker` is a command-line service for running arbitrary Linux processes.

To use it, first run the server:

```sh
./server
```

Then, in a separate terminal instance, set environment variables for the pre-determined username and password:

```sh
$ export username="default_user"
$ export pw="123456"
```

In that same terminal instance, you may now begin scheduling jobs for the server to execute, for example:

```sh
$ ./worker run echo hello
Job scheduled. JobID: 1
$ ./worker run echo world
Job scheduled. JobID: 2
$ ./worker out 1
hello
$ ./worker out 2
world
```

To view the usage for additional commands, run `./worker help`
