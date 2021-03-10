# k8s-webhooks-stuffs

## This project docker image

- For the images of this project you can see this dockerhub [repo](https://hub.docker.com/repository/docker/shahincsejnu/mutator-webhook)
- use the latest image (image >= v1.0.5)

## Setup this Mutating Webhook

```
kubectl apply -f mutating-webhook-configuration.yaml
kubectl apply -f mutating-webhook-deployment.yaml
kubectl apply -f mutating-webhook-service.yaml
```

## Set Up the CA Certificate

- Once the webhook runs (give it a few seconds to initialize), the CA certificate can be downloaded by executing a curl command within the container. To retrieve the base64 encoded version of this ca.pem, use the following command:

    ```
    kubectl exec -it pods <respective_pod_name_deployed_by_mutating-webhook-deployment> sh 
    wget -q -O- http://127.0.0.1:8080/ca.pem?base64
    ```

- The output of this command should replace the base64 string in caBundle in mutating-webhook-configuration.yaml:

    ```
    caBundle: "<pre_string>"  # <= replace this string within quotes
    ```

- Then re-apply the webhook again:

    ```
    kubectl apply -f mutating-webhook-configuration.yaml
    ```
  
- Now apply the teployment object:
    ```
    kc apply -f teployment-obj.yaml
    ```