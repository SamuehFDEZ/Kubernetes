# Cheatsheet de kubernetes


## Crear pods

## Start a nginx pod
```bash
kubectl run nginx --image=nginx
```
## Start a hazelcast pod and let the container expose port 5701
```bash
kubectl run hazelcast --image=hazelcast/hazelcast --port=5701
```
## Start a hazelcast pod and set environment variables "DNS_DOMAIN=cluster" and "POD_NAMESPACE=default" in the container
```bash
kubectl run hazelcast --image=hazelcast/hazelcast --env="DNS_DOMAIN=cluster" --env="POD_NAMESPACE=default"
```
## Start a hazelcast pod and set labels "app=hazelcast" and "env=prod" in the container
```bash
kubectl run hazelcast --image=hazelcast/hazelcast --labels="app=hazelcast,env=prod"
```
## Dry run; print the corresponding API objects without creating them
```bash
kubectl run nginx --image=nginx --dry-run=client
```
## Start a nginx pod, but overload the spec with a partial set of values parsed from JSON
```bash
kubectl run nginx --image=nginx --overrides='{ "apiVersion": "v1", "spec": { ... } }'
```
## Start a busybox pod and keep it in the foreground, don't restart it if it exits
```bash
kubectl run -i -t busybox --image=busybox --restart=Never
```
## Start the nginx pod using the default command, but use custom arguments (arg1 .. argN) for that command
```bash
kubectl run nginx --image=nginx -- <arg1> <arg2> ... <argN>
```
## Start the nginx pod using a different command and custom arguments
```bash
kubectl run nginx --image=nginx --command -- <cmd> <arg1> ... <argN>
```
--- 

## Obtener pods

```bash
kubectl get pod
```
--- 

## Borrar pods

```bash
kubectl delete pod <nombre>
```

## Obtener pod detalladamente
```bash
kubectl get pod <nombre> -o yaml
```


## Mapear puerto en kubernetes
```bash
kubectl port-forward pod/podtest 8080:80
```

## Entrar en la linea de comandos del pod
```bash
kubectl exec -it podtest -- sh
```

## Ver logs de pods
```bash
kubectl logs podtest -f
```

## Manifiestos de kubernetes
#### los yamls 

Ver template de pod en
https://kubernetes.io/docs/concepts/workloads/pods/


```
apiVersion: v1
kind: Pod
metadata:
  name: podtest2
spec:
  containers:
  - name: cont1
    image: nginx:alpine
    # The pod template ends here
```



## Ver versiones y recursos de api de kubernetes
```bash
kubectl api-versions
```

```bash
kubectl api-resources | grep Pod
```


## Labels
Dentro de metadata se asigna un array de labels, ejemplo:
Siempre asignar al menos el label de app

```bash
metadata:
  name: podtest2
  labels:
    app: front
    env: dev
```


```bash
kubectl get pods -l app=backend
```


## Problemas de los pods

- Sin self-healing (No se autoregeneran)
- Si quieres crear masivamente pods tienes que hacerlo manualmente en el yaml
- Los pods no tienen autorefresh, es decir, no se actualizan solos


## Replicaset

- Objeto superior a los pods
- Se "adueña" de ellos y los crea
- Agrega al metadata de los pods el valor owner, referenciando a qué replicaset pertenecen
- Otro replicaset no puede tomar un pod que ya tenga owner

### Para obtener replicasets por shortname
```bash
kubectl get rs
```
rs de replicaset eso se puede consultar en

```bash
kubectl api-resources
```