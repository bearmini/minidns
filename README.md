minidns
=======

A tiny DNS server with very limited features and the simplest configuration.

## how to install

```
go get github.com/bearmini/minidns
```


## how to run

```
minidns --config config_example.com
```

in another terminal:
```
dig example.com @localhost -p 8053

; <<>> DiG 9.8.3-P1 <<>> example.com @localhost -p 8053
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 14292
;; flags: qr aa rd; QUERY: 1, ANSWER: 2, AUTHORITY: 0, ADDITIONAL: 0
;; WARNING: recursion requested but not available

;; QUESTION SECTION:
;example.com.			IN	A

;; ANSWER SECTION:
example.com.		60	IN	A	8.8.8.8
example.com.		60	IN	A	8.8.4.4

;; Query time: 2 msec
;; SERVER: ::1#8053(::1)
;; WHEN: Tue Apr 18 15:55:17 2017
;; MSG SIZE  rcvd: 83

```
