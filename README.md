# Worker

`worker` is a command-line service for running arbitrary Linux processes.

To use it, first make the builds, which will be placed in a new `bin/` directory:

```sh
$ make build
```

Now you may run the server:

```sh
$ export port="8080" # Optionally set port (default 443). Ports <= 1024 require the server to be started using 'sudo'.
$ bin/server
Listening...
```

Then, in a separate terminal instance, set environment variables for the pre-determined username and password:

```sh
$ export port="8080" # If the port was set for the server, it must be the same for the client.
$ export username="default_user"
$ export pw="123456"
```

In that same terminal instance, you may now begin scheduling jobs for the server to execute, for example:

```sh
$ bin/worker run echo hello
Ht9piRvJVMWq5CnTShXMkY # After a job is scheduled, its id is printed.
$ bin/worker run echo world
nya8Z45ei5BTkgWdqN3NWc
$ bin/worker out Ht9piRvJVMWq5CnTShXMkY # Use a job id to query output.
hello
$ bin/worker out nya8Z45ei5BTkgWdqN3NWc
world
```

To view the usage for additional commands, run `bin/worker help`
