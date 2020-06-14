# Monitoring The Near Node with Prometheus and Grafana 

In this guide, you will learn how to setup Prometheus node exporter and [near exporter](https://github.com/masknetgoal634/near-prometheus-exporter) on a Near node to export metrics to the Prometheus server and monitor them with Grafana.

## Run Node Exporter on the Node

First we need to deploy [near prometheus exporter](https://github.com/masknetgoal634/near-prometheus-exporter) service to collect custom metrics from the near node using json-rpc.

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

    git clone https://github.com/masknetgoal634/near-prometheus-exporter

    cd near-prometheus-exporter

    sudo docker build -t near-prometheus-exporter .

```
sudo docker run -dit \
    --restart always \
    --name near-exporter \
    --network=host \
    -p 9333:9333 \
    near-prometheus-exporter:latest /dist/main -accountId <YOUR_POOL_ID>
```

Open 9333 port in your server firewall as Prometheus reads metrics on this port.

## Collect Metrics From The Near Node

Open 3030 port in your server firewall to collect metrics from the near node itself

## Configure the Near Node as Target on Prometheus Server

Now that we have the node and near exporters up and running on the node, we have to add targets on the Prometheus server configuration.

>You can run Prometheus server on your home computer or even on Raspberry Pi 

    git clone https://github.com/masknetgoal634/near-prometheus-exporter

    cd near-prometheus-exporter/etc

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

  - job_name: near-node
    scrape_interval: 15s
    static_configs:
    - targets: ['<NODE_IP_ADDRESS>:3030']
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

### Run Grafana

```
sudo chown -R 472:472 grafana

sudo docker run -dit \
    --restart always \
    --volume $(pwd)/grafana:/var/lib/grafana \
    --volume $(pwd)/grafana/provisioning:/etc/grafana/provisioning \
    --volume $(pwd)/grafana/custom.ini:/etc/grafana/grafana.ini \
    --user 472 \
    --network=host \
    --name grafana \
    -p 3000:3000 grafana/grafana
```

Open in your favorite browser `http://localhost:3000`

![](https://raw.githubusercontent.com/masknetgoal634/near-prometheus-exporter/master/guide/img/image0.png)

Username: admin
Password: admin

Open `Near Node Exporter Full` dashboard

![](https://raw.githubusercontent.com/masknetgoal634/near-prometheus-exporter/master/guide/img/image1.png)

![](https://raw.githubusercontent.com/masknetgoal634/near-prometheus-exporter/master/guide/img/image2.png)

Set up refresh time at the top-right corner of the dashboard:
![](https://raw.githubusercontent.com/masknetgoal634/near-prometheus-exporter/master/guide/img/refresh_time.png)

### Grafana Email Notification Alert

First we need to edit the grafana config file `etc/grafana/custom.ini`

```
[smtp]
enabled = true
host = smtp.gmail.com:587 
user = <your_gmail_address>
password = <your_gmail_password>
;cert_file =
;key_file =
skip_verify = true
from_address = <your_gmail_address>
from_name = Grafana
```

Now we need to configure an Email Alert channel in the alerting serction:

![](https://raw.githubusercontent.com/masknetgoal634/near-prometheus-exporter/master/guide/img/email_channel.png)

Enter your an email.

Click on "Send Test" button.
If you getting error you need to access the grafana app the following url and accept the app:

 https://accounts.google.com/DisplayUnlockCaptcha

Now click again on "Send Test" button and look into your email account inbox.

Finally we are ready to create our first alert!

Open "Near Node Exporter Full" dashboard end edit "Peer Connections Total" panel:  

![](https://raw.githubusercontent.com/masknetgoal634/near-prometheus-exporter/master/guide/img/email_alert.png)

You need to set a threshold for the peer connections, for instance you may set 5.
When the peer connections will be below 5 you will see an alert in your email inbox.  

Enjoy

Updated: 13.06.2020
