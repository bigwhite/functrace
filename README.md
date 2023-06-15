# functrace

## Quick Start

### Install gen

```
$go get github.com/pengxuan37/functrace/cmd/gen
```
And make sure the gen executable file is configured as part of the PATH env variable.

### Generate trace entry for your packages

See the demo case: https://github.com/pengxuan37/functrace/tree/main/examples/gen-demo

### Add trace in batches

You can use the scripts/batch_add_trace.sh to add trace in batches for all go source files in some repo.

for example:

let's add trace in batches for github.com/panjf2000/gnet, the steps is below:

- git clone https://github.com/panjf2000/gnet.git
- cd gnet
- execute the following command:

```
bash -c "$(curl -fsSL https://raw.githubusercontent.com/bigwhite/functrace/main/scripts/batch_add_trace.sh)"
```

>If the access to github.com is very slow in your area, then you can copy scripts/batch_add_trace.sh of functrace to gnet dir and execute ```bash batch_add_trace.sh```.

The command will output:

```
[gen -w ./ringbuffer/ring_buffer_test.go]
add trace for ./ringbuffer/ring_buffer_test.go ok
[gen -w ./ringbuffer/ring_buffer.go]
add trace for ./ringbuffer/ring_buffer.go ok
... ...
[gen -w ./internal/netpoll/queue/queue.go]
no trace added for ./internal/netpoll/queue/queue.go
[gen -w ./gnet.go]
add trace for ./gnet.go ok
[gen -w ./acceptor_windows.go]
add trace for ./acceptor_windows.go ok
```

After debugging, you can use 'git checkout .' to recover the origin status of your local workspace.
