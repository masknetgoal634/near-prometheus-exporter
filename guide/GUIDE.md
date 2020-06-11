# Monitor Near Node using Prometheus and Grafana

In this guide, you will learn how to setup Prometheus node exporter and [near exporter](https://github.com/masknetgoal634/near-prometheus-exporter) on a Near node to export metrics to the Prometheus server and monitor them with Grafana.

## Run Node Exporter on the Node

```
sudo docker run -dit \
    --restart always \
    --volume /proc:/host/proc:ro \
    --volume /sys:/host/sys:ro \
    --volume /:/rootfs:ro \
    --name node-exporter \
    -p 9100:9100 prom/node-exporter:latest \
    --path.procfs=/host/proc \
    --path.sysfs=/host/sys
```

Open 9100 port in your server firewall as Prometheus reads metrics on this port.

## Run Near Exporter on the Node

`git clone https://github.com/masknetgoal634/near-prometheus-exporter`

`cd near-prometheus-exporter`

`sudo docker build -t near-prometheus-exporter .`

```
sudo docker run -dit \
    --restart always \
    --name near-exporter \
    -p 9333:9333 \
    masknetgoal634/near-prometheus-exporter:latest /dist/main -accountId <YOUR_POOL_ID>
```

Open 9333 port in your server firewall as Prometheus reads metrics on this port.

## Configure the Near Node as Target on Prometheus Server

Now that we have the node and near exporters up and running on the node, we have to add targets on the Prometheus server configuration.

>You can run Prometheus server on your home computer or even on Raspberry Pi 

`git clone https://github.com/masknetgoal634/near-prometheus-exporter`

`cd near-prometheus-exporter/etc`

open `prometheus/prometheus.yml` and add an ip address of your node:

```
  ...
  - job_name: node
    scrape_interval: 5s
    static_configs:
    - targets: ['<NODE_IP_ADDRESS>:9100']

  - job_name: near-exporter
    scrape_interval: 15s
    static_configs:
    - targets: ['<NODE_IP_ADDRESS>:9333']
  ...
```

## Run Prometheus on your server monitoring machine

```
sudo docker run -dti \
    --restart always \
    --volume $(pwd)/prometheus:/etc/prometheus/ \
    --name prometheus \
    --network=host \
    -p 9090:9090 prom/prometheus:latest \
    --config.file=/etc/prometheus/prometheus.yml
```

#### Run Grafana

```
sudo chown -R 472:472 grafana
sudo docker run -dit \
    --restart always \
    --volume $(pwd)/grafana:/var/lib/grafana \
    --volume $(pwd)/grafana/provisioning:/etc/grafana/provisioning \
    --user 0 \
    --name grafana \
    -p 3000:3000 grafana/grafana
```

Open in your favorite browser `http://localhost:3000`

![](https://raw.githubusercontent.com/masknetgoal634/near-prometheus-exporter/master/guide/img/image0.png)

login "admin"
password "admin"

Open `Near Node Exporter Full` dashboard

![](https://raw.githubusercontent.com/masknetgoal634/near-prometheus-exporter/blob/master/guide/img/image1.png)

![](https://raw.githubusercontent.com/masknetgoal634/near-prometheus-exporter/blob/master/guide/img/image2.png)

### Configure Alert in Grafana

Add a telegram channel with your Telegram TokenId and ChannelId

![](https://raw.githubusercontent.com/masknetgoal634/near-prometheus-exporter/blob/master/guide/img/image3.png)

Also add an alert when blocks did not produced during 5 minutes

![](https://raw.githubusercontent.com/masknetgoal634/near-prometheus-exporter/blob/master/guide/img/image4.png)

Thats it!
