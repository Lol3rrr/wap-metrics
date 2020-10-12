# wap-metrics
A simple helper to export wifi AP metrics for raspberry pi

## Usage
### Basic
To convert the raw `iw {interface} station dump` into prometheus metrics run

`iw {interface} station dump | ./wap-metrics`

This will then output the metrics in the right format for prometheus to stdout

### Automatic - Node exporter
#### Requirements
 * Node-Exporter running on the raspberry pi
 * Textfile collector enabled for the node-exporter

#### Setup
To automatically collect these metrics, the following should be run on a automatically using something like cron.

`/sbin/iw {interface} station dump | /absolute/path/wap-metrics | sponge /abolute/path/metrics/folder/wap.prom`

This will collect the latest metrics about the AP and then use sponge to write them to a wap.prom file where node-exporter can find it.
