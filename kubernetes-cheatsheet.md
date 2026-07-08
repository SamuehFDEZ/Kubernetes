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

## Para agregar labels a pods sin owner
```bash
kubectl label pods podtest1 app=pod-label
```

Al crear dos pods diferentes con el mismo label, el replicaSet adopta esos pods como suyos y agrega a la metadata ownerReferences la referencia del replicaSet

¡Problema! los pods son totalmente distintos pero para el replicaset son iguales debido al label que comparten

Es por ello que los pods han de ser creados por unidades u objetos mayores, ya sean replicasets o deployments


## Problemas de ReplicaSet

El concepto general es que debe mantener un numero n de replicas de un pod segun lo que digamos en el manifiesto (yaml)

En el caso de que se quiera hacer un cambio en plano al pod (directamente a él) no ocurrirá nada, ya que el replicaset solo mira el numero de pods que corresponden al label definido dentro de la etiqueta metadata, por lo que no puede cambiar los pods ni configuraciones


## Deployments

Un deployment es un objeto que está por encima de un replicaset y por encima del pod

MaxAvailable -> cuantos pods voy a permitir que esté fuera de servicio, por defecto 25%

MaxSearch -> cuanto voy a permitir adicional al 100% para que se creen pods nuevos, es decir si ya tengo 4 pods, cuantos permito tener de más

Kubernetes por defecto mantiene 10 replicaSets

Para mostrar labels de un deployment por ejemplo
```bash
kubectl get deployment --show-labels
```

Comando para verificar si el rollout del deployment ha sido un exito

```bash
kubectl rollout status deployment <nombreDeployment>
kubectl rollout status deployment deployment-test
```

### OwnerReferences en deployment

Un pod va a tener como ownerReference un replicaset y un replicaset va a tener como ownerReference a un deployment siempre

No pueden saltarse, es decir, un pod no puede tener como ownerReference a un deployment


## Rolling Updates de deployments
```bash
kubectl apply -f deployment.yaml
```

Si aplicamos el comando de apply -f al yaml del deployment con algun cambio éste, dependiendo de lo que tenga configurado, eliminará y creará pods con las nuevas especificaciones

Con el comando 

```bash
kubectl rollout status deployment <nombreDeployment>
```

Vemos si ha sido exitoso el rollout, es decir, el cambio aplicado a los pods

Si hacemos un describe

```bash
kubectl describe deploy deployment-test
```

Podemos ver el siguiente output

```
Events:
  Type    Reason             Age                From                   Message
  ----    ------             ----               ----                   -------
  Normal  ScalingReplicaSet  43h                deployment-controller  Scaled up replica set deployment-test-6cf85c55cf to 3
  Normal  ScalingReplicaSet  2m36s              deployment-controller  Scaled up replica set deployment-test-69b6fb5cb6 to 1
  Normal  ScalingReplicaSet  2m31s              deployment-controller  Scaled down replica set deployment-test-6cf85c55cf to 2 from 3
  Normal  ScalingReplicaSet  2m31s              deployment-controller  Scaled up replica set deployment-test-69b6fb5cb6 to 2 from 1
  Normal  ScalingReplicaSet  2m26s              deployment-controller  Scaled down replica set deployment-test-6cf85c55cf to 1 from 2
  Normal  ScalingReplicaSet  2m25s              deployment-controller  Scaled up replica set deployment-test-69b6fb5cb6 to 3 from 2
  Normal  ScalingReplicaSet  2m19s              deployment-controller  Scaled down replica set deployment-test-6cf85c55cf to 0 from 1
  Normal  ScalingReplicaSet  83s                deployment-controller  Scaled up replica set deployment-test-6cf85c55cf to 1 from 0
  Normal  ScalingReplicaSet  81s                deployment-controller  Scaled down replica set deployment-test-69b6fb5cb6 to 2 from 3
  Normal  ScalingReplicaSet  80s                deployment-controller  Scaled up replica set deployment-test-6cf85c55cf to 2 from 1
  Normal  ScalingReplicaSet  74s (x3 over 77s)  deployment-controller  (combined from similar events): Scaled down replica set deployment-test-69b6fb5cb6 to 0 from 1
```