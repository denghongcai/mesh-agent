#!/bin/bash

ETCD_HOST=$(ip addr show docker0 | grep 'inet\b' | awk '{print $2}' | cut -d '/' -f 1)
ETCD_PORT=2379
ETCD_URL=http://$ETCD_HOST:$ETCD_PORT

echo ETCD_URL = $ETCD_URL

if [[ "$1" == "consumer" ]]; then
  echo "Starting consumer agent..."
  /root/dists/agent --etcdEndPoint $ETCD_URL --role consumer --listenPort 20000 > /root/logs/std.log 2>&1 
elif [[ "$1" == "provider-small" ]]; then
  echo "Starting small provider agent..."
  /root/dists/agent --etcdEndPoint $ETCD_URL --weight 1 --role provider --interfaceName com.alibaba.dubbo.performance.demo.provider.IHelloService --listenPort 30000 --servicePort 20889 > /root/logs/std.log 2>&1 
elif [[ "$1" == "provider-medium" ]]; then
  echo "Starting medium provider agent..."
  /root/dists/agent --etcdEndPoint $ETCD_URL --weight 2 --role provider --interfaceName com.alibaba.dubbo.performance.demo.provider.IHelloService --listenPort 30001 --servicePort 20890 > /root/logs/std.log 2>&1 
elif [[ "$1" == "provider-large" ]]; then
  echo "Starting large provider agent..."
  /root/dists/agent --etcdEndPoint $ETCD_URL --weight 3 --role provider --interfaceName com.alibaba.dubbo.performance.demo.provider.IHelloService --listenPort 30002 --servicePort 20891 > /root/logs/std.log 2>&1 
else
  echo "Unrecognized arguments, exit."
  exit 1
fi
