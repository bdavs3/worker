# Worker

`worker` is a command-line service for running arbitrary Linux processes. The Linux binaries can be found in the `bin` directory.

To use it, first run the server:

```sh
$ ./server # If the server does not persist, run it using 'sudo'.
Listening...
```

Then, in a separate terminal instance, set environment variables for the pre-determined username and password:

```sh
$ export username="default_user"
$ export pw="123456"
```

In that same terminal instance, you may now begin scheduling jobs for the server to execute, for example:

```sh
$ ./worker run echo hello
Ht9piRvJVMWq5CnTShXMkY # After a job is scheduled, its id is printed.
$ ./worker run echo world
nya8Z45ei5BTkgWdqN3NWc
$ ./worker out Ht9piRvJVMWq5CnTShXMkY # Use a job id to query output.
hello
$ ./worker out nya8Z45ei5BTkgWdqN3NWc
world
```

To view the usage for additional commands, run `./worker help`
