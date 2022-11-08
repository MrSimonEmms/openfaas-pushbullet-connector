# OpenFaaS Pushbullet Connector Chart

## Installation


Your first job is to get a Pushbullet API token (see [main docs](../../)). Once acquired, you will to create a
secret called `pushbullet`.

```shell
export PUSHBULLET_TOKEN=xxx

kubectl create secret generic \
		-n openfaas \
		pushbullet \
		--from-literal=token=${PUSHBULLET_TOKEN} \
		--dry-run=client -o yaml | \
		kubectl replace --force -f -
```

Now that's installed, you will need to install the connector with Helm.

```shell
helm upgrade \
  --atomic \
  --cleanup-on-fail \
  --create-namespace \
  --install \
  --namespace="openfaas" \
  --repo="https://helm.simonemms.com" \
  --reset-values \
  --wait \
  pushbullet-connector \
  openfaas-pushbullet-connector
```

You can watch the Connector logs to see it invoke your functions:

```shell
kubectl logs -f -n openfaas deploy/pushbullet-connector
```
