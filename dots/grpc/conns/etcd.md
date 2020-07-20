# etcd server
the etcd 3.4.9 on default port 2379 , (2380 is for peer)

[download](https://github.com/etcd-io/etcd/releases)

## cluster
[see](https://github.com/etcd-io/etcd/blob/master/Documentation/op-guide/clustering.md)

static
On each machine, start etcd with these flags:

```
$ etcd --name etcdnode1 --initial-advertise-peer-urls http://10.0.1.10:2380 \
  --listen-peer-urls http://10.0.1.10:2380 \
  --listen-client-urls http://10.0.1.10:2379,http://127.0.0.1:2379 \
  --advertise-client-urls http://10.0.1.10:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster etcdnode1=http://10.0.1.10:2380,etcdnode2=http://10.0.1.11:2380,etcdnode3=http://10.0.1.12:2380 \
  --initial-cluster-state new
```
```
$ etcd --name etcdnode2 --initial-advertise-peer-urls http://10.0.1.11:2380 \
  --listen-peer-urls http://10.0.1.11:2380 \
  --listen-client-urls http://10.0.1.11:2379,http://127.0.0.1:2379 \
  --advertise-client-urls http://10.0.1.11:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster etcdnode1=http://10.0.1.10:2380,etcdnode2=http://10.0.1.11:2380,etcdnode3=http://10.0.1.12:2380 \
  --initial-cluster-state new
```
```
$ etcd --name etcdnode3 --initial-advertise-peer-urls http://10.0.1.12:2380 \
  --listen-peer-urls http://10.0.1.12:2380 \
  --listen-client-urls http://10.0.1.12:2379,http://127.0.0.1:2379 \
  --advertise-client-urls http://10.0.1.12:2379 \
  --initial-cluster-token etcd-cluster-1 \
  --initial-cluster etcdnode1=http://10.0.1.10:2380,etcdnode2=http://10.0.1.11:2380,etcdnode3=http://10.0.1.12:2380 \
  --initial-cluster-state new
```

# etcd browser
1. [etcdkeeper](https://github.com/evildecay/etcdkeeper)
2. [e3w](https://github.com/xiaowei520/e3w)
3. [etcd-manager](https://www.npmjs.com/package/etcd-manager)
