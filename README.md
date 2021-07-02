# Deployment Editor
This is an example project, is it not meant to be used in production or staging. Should only be used on a local cluster.

# What this does
This project simply takes advantage of kubernetes mutating webhooks. In the simplest form when a deployment manifest is applied to the cluster the mutating webhook will receive the manifest before it's applied, do any changes and return to the cluster which is then applied.  

## Annotations
This controller makes use of manifest annotations, in this example it looks for an annotation called: `jackfazackerley.com/should-edit-replicas` the value should be either `"true"` or `"false"`. If the annotation isn't present then nothing will be done, it's the same as the annotation value being `"false"`. If the annotation is present and is `"true"` it will set the replicas to the value of `--replicas` is. [manager.yaml](config/manager/manager.yaml)

# How to use
In order to deploy to k8s first two binary files need to be downloaded.
```shell
make controller-gen
make kustomize
```

This requires cert-manager to be installed for CA auth from the cluster. To install cert-manager run the following:
```shell
helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.4.0 \
  --set installCRDs=true
```

To install the controller to the cluster simply run:
```shell
make deploy
```

and to uninstall run:
```shell
make undeploy
helm delete cert-manager
```

To deploy an example deployment use the following: 
```shell
kubectl apply -f config/samples/deployment.yaml
```

To delete the example deployment use the following:
```shell
kubectl delete -f config/samples/deployment.yaml
```