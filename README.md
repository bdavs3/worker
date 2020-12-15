# Worker

Worker is a command-line service for running arbitrary Linux processes.

To use it, first run the server:

```sh
./server
```

Then, in a separate terminal instance, begin scheduling jobs for the server to execute, for example:

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
