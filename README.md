# Worker

`worker` is a command-line service for running arbitrary Linux processes.

I completed this project as the final stage of the hiring process for a backend engineering role at [Teleport](https://goteleport.com). The [challenge spec](https://github.com/bdavs3/worker/blob/master/Challenge.pdf) is in the repository if you'd like more details. Before this project, I had zero Golang programming experience and my knowledge of backend was limited. I was ultimately successfully in fulfilling the L1 requirements, but was unfortunately not offered the job. I did, however, learn that I really like programming in Go and building backend systems.

As I am now going to continue pursuing similar employment opportunities, this repository serves to demonstrate the progress I made in my first month of working in this subfield of computer science. Since my interview team at Teleport was constantly in touch to provide direction and identify shortcomings in my code, you can also check out the closed [pull requests](https://github.com/bdavs3/worker/pulls?q=is%3Apr+is%3Aclosed) to understand how I approach problem-solving in a remote, open-source environment.

To use the service, first make the builds, which will be placed in a new `bin` directory. Make sure you build for your operating system:

```sh
$ make windows
$ make mac
$ make linux
```

Once the builds are complete, `cd` into the `bin` directory.

Now you may run the server:

```sh
$ export port="8080" # Optionally set port (default 443). Ports <= 1024 require the server to be started using 'sudo'.
$ ./server
Listening...
```

Then, in a separate terminal instance, `cd` to the `bin` directory, and set environment variables for the pre-determined username and password:

```sh
$ export port="8080" # If the port was set for the server, it must be the same for the client.
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
