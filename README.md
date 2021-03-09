# k8s-webhooks-stuffs

## Setup this Mutating Webhook

```
kubectl apply -f mutating-webhook-configuration.yaml
kubectl apply -f mutating-webhook-deployment.yaml
kubectl apply -f mutating-webhook-service.yaml
```

## Set Up the CA Certificate

Once the webhook runs (give it a few seconds to initialize), the CA certificate can be downloaded by executing a curl command within the container. To retrieve the base64 encoded version of this ca.pem, use the following command:

```
kubectl exec -it $(kubectl get pods --no-headers -o custom-columns=":metadata.name") -- wget -q -O- localhost:8080/ca.pem?base64
```

The output of this command should replace the base64 string in caBundle in mutating-webhook-configuration.yaml:

```
    caBundle: "cGxhY2Vob2xkZXIK" # <= replace this string within quotes
```

Then re-apply the webhook again:

```
kubectl apply -f mutating-webhook-configuration.yaml
```