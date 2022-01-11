# gmqtt-proxy-plugin
It's a plugin created for [Gmqtt broker](https://github.com/DrmagicE/gmqtt), to make it act like a mqtt proxy, with messages buffering support when the destination broker is down.

## How to install

1. Clone [Gmqtt project](https://github.com/DrmagicE/gmqtt)

2. Edit `plugin_imports.yml`file in the root directory and add proxy plugin:

```yml
packages:
  - admin
  - prometheus
  - federation
  - auth 
  # add proxy plugin in the end of the list
  - proxy 
```

3. Edit `cmd/gmqttd/default_config.yml` file and add proxy plugin:

```yml
plugin_order:  
  - prometheus
  - admin
  - federation
  # add proxy plugin here
  - proxy
```

4. Go to `plugin` directory and clone this project:

```shell
git clone https://github.com/Oliveirakun/gmqtt-proxy-plugin.git
```

5. Rename the directory:

```shell
mv gmqtt-proxy-plugin proxy
```

5. Go back to root directory and install the plugin dependencies:

```shell
go get github.com/eclipse/paho.mqtt.golang
```

6. Export an environment variable with the destination broker uri:

```shell
export REMOTE_MQTT_BROKER=tcp://127.0.0.1:1884
```

7. Run the project:

```shell
make run
```
## Routes

Route | Description
------------ | -------------
**/v1/toogle?stop=true** | Stop sending messages to destination broker
**/v1/toogle?stop=false** | Restart sending messages to destination broker
**/v1/toogle?stop=false&broker-uri=tcp://127.0.0.1** | Restart sending messages to a different broker



