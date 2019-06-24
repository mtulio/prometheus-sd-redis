# prometheus-sd-redis

Prometheus File Service Discovery for Redis Cluster and Standalone.

This job discovery redis endpoints and use [redis exporter](https://github.com/oliver006/redis_exporter#prometheus-configuration-to-scrape-multiple-redis-hosts) to scrape metrics from multiple redis hosts.

## Config

First of all you need to create an config file for each service to discover. In general each server could be an group of common servers to be disvered, E.g: an Redis Cluster, then one extra labels could be added that will be attached on the list of labels of the jobs.

See an example of config:

```json
{
    "services": [
        {
            "type": "cluster",
            "job": "redis",
            "url": "redis-cluster-ecache.sample:6379",
            "labels":
                {
                    "type": "cluster",
                    "name": "rcluster-elasticache"
                }
        },
        {
            "type": "cluster",
            "job": "redis",
            "url": "redis-cluster-ec2.sample:7500",
            "labels":
                {
                    "type": "cluster",
                    "name": "rcluster-ec2"
                }
        }
    ]
}
```

* `services`: list of services
* `type`: could be `cluster` or `standalone`, is the type of the endpoint, and the way that the discovery will be done.
* `job`: the job name of Prometheus.
* `url`: the URL of one node of the server. Baseline to discovery
* `labels`: extra labels that will be attached to the discovery file used by Prometheus' SD `file_sd_configs`.

## Usage

* One shot running:

```bash
./prometheus-sd-redis --in-file ./contrib/config-sample.json --out-file ./contrib/out.json
```

* Data generated on `./contrib/out.json`

```json
[
 {
  "targets": [
   "172.18.92.63:6379",
   "172.18.2.23:6379",
   "172.18.83.251:6379",
   "172.18.3.186:6379",
   "172.18.89.24:6379",
   "172.18.15.69:6379",
   "172.18.82.197:6379",
   "172.18.1.165:6379"
  ],
  "labels": {
   "job": "redis",
   "name": "rcluster-elasticache",
   "type": "cluster"
  }
 },
 {
  "targets": [
   "10.250.111.80:8800",
   "10.250.110.120:8801",
   "10.250.108.246:8800",
   "10.250.108.47:8800",
   "10.250.108.223:8801",
   "10.250.110.238:8801",
   "10.250.109.131:8801",
   "10.250.108.246:8801",
   "10.250.110.6:8801",
   "10.250.108.143:8801",
   "10.250.110.27:8801",
   "10.250.109.106:8800",
   "10.250.108.86:8801",
   "10.250.111.21:8801",
   "10.250.109.130:8800",
   "10.250.109.9:8800",
   "10.250.110.99:8801",
   "10.250.109.131:8800",
   "10.250.108.86:8800",
   "10.250.111.73:8800",
   "10.250.110.238:8800",
   "10.250.110.209:8800",
   "10.250.110.106:8801",
   "10.250.109.130:8801",
   "10.250.110.99:8800",
   "10.250.111.73:8801",
   "10.250.110.27:8800",
   "10.250.110.106:8800",
   "10.250.108.47:8801",
   "10.250.110.156:8801",
   "10.250.108.223:8800",
   "10.250.111.21:8800",
   "10.250.110.209:8801",
   "10.250.111.80:8801",
   "10.250.110.156:8800",
   "10.250.109.9:8801",
   "10.250.110.6:8800",
   "10.250.109.106:8801",
   "10.250.110.120:8800",
   "10.250.108.143:8800"
  ],
  "labels": {
   "job": "redis",
   "name": "rcluster-ec2",
   "type": "cluster"
  }
 }
]
```

## Contribute

Open and PR or Issue. =)