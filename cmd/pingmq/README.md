pingmq
======

The following commands will run pingmq as a server, pinging the 8.8.8.0/28 CIDR block, and publishing the results to /ping/success/{ip} and /ping/failure/{ip} topics every 30 seconds. `sudo` is needed because we are using RAW sockets and that requires root privilege.

```
$ go build
$ sudo ./pingmq server -p 8.8.8.0/28 -i 30
```

The following command will run pingmq as a client, subscribing to /ping/failure/+ topic and receiving any failed ping attempts.

```
$ ./pingmq client -t /ping/failure/+
8.8.8.6: Request timed out for seq 1
```

The following command will run pingmq as a client, subscribing to /ping/failure/+ topic and receiving any failed ping attempts.

```
$ ./pingmq client -t /ping/success/+
8 bytes from 8.8.8.8: seq=1 ttl=56 tos=32 time=21.753711ms
```

One can also subscribe to a specific IP by using the following command.

```
$ ./pingmq client -t /ping/+/8.8.8.8
8 bytes from 8.8.8.8: seq=1 ttl=56 tos=32 time=21.753711ms
```