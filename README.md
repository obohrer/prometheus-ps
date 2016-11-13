# prometheus-ps
Export the state of your processes to prometheus

## build
Install go & glide

Run ```make```

## Running the exporter
The exporter will load the conf.json file specified by : ```--config```

The config file specifies a watchlist (process names to export) and the port to use to export the ps.

Example of config :

```json
{"WatchList":["top", "usbagent", "mdworker", "iTerm2"],
 "Port":9105}
```
