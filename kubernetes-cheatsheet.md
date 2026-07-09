<div align="center">

# ☸️ Kubernetes Cheatsheet
 
### Guía de referencia rápida: Pods, ReplicaSets, Deployments y Servicios
 
---
 
📦 `kubectl` · 🐳 `Docker Desktop` · 🧭 `Manifiestos YAML`
 
</div>

---
 
## 📑 Índice
 
1. [🚀 Crear pods](#-crear-pods)
2. [🔍 Obtener y gestionar pods](#-obtener-y-gestionar-pods)
3. [📄 Manifiestos de Kubernetes](#-manifiestos-de-kubernetes)
4. [🏷️ Labels](#️-labels)
5. [⚠️ Problemas de los pods](#️-problemas-de-los-pods)
6. [🧬 ReplicaSet](#-replicaset)
7. [🚢 Deployments](#-deployments)
8. [🔄 Rolling Updates](#-rolling-updates-de-deployments)
9. [📜 Histórico y revisiones](#-histórico-y-revisiones-de-un-deployment)
10. [⏪ Rollbacks](#-rollbacks)
11. [🌐 Servicios](#-servicios)
---
 
## 🚀 Crear pods
 
### ▶️ Iniciar un pod nginx
```bash
kubectl run nginx --image=nginx
```
 
### ▶️ Iniciar un pod hazelcast exponiendo el puerto 5701
```bash
kubectl run hazelcast --image=hazelcast/hazelcast --port=5701
```
 
### ▶️ Iniciar un pod hazelcast con variables de entorno
Define `DNS_DOMAIN=cluster` y `POD_NAMESPACE=default` en el contenedor:
```bash
kubectl run hazelcast --image=hazelcast/hazelcast --env="DNS_DOMAIN=cluster" --env="POD_NAMESPACE=default"
```
 
### 🏷️ Iniciar un pod hazelcast con labels
Define `app=hazelcast` y `env=prod`:
```bash
kubectl run hazelcast --image=hazelcast/hazelcast --labels="app=hazelcast,env=prod"
```
 
### 🧪 Dry run
Imprime los objetos de la API sin crearlos:
```bash
kubectl run nginx --image=nginx --dry-run=client
```
 
### ⚙️ Pod con spec sobrescrita mediante JSON
```bash
kubectl run nginx --image=nginx --overrides='{ "apiVersion": "v1", "spec": { ... } }'
```
 
### 📦 Pod busybox interactivo, sin reinicio
```bash
kubectl run -i -t busybox --image=busybox --restart=Never
```
 
### 🎯 Comando por defecto con argumentos personalizados
```bash
kubectl run nginx --image=nginx -- <arg1> <arg2> ... <argN>
```
 
### 🛠️ Comando y argumentos personalizados
```bash
kubectl run nginx --image=nginx --command -- <cmd> <arg1> ... <argN>
```
 
---
 
## 🔍 Obtener y gestionar pods
 
### 📋 Obtener pods
```bash
kubectl get pod
```
 
### 🗑️ Borrar pods
```bash
kubectl delete pod <nombre>
```
 
### 🔬 Obtener detalle de un pod
```bash
kubectl get pod <nombre> -o yaml
```
 
### 🔌 Mapear puerto en Kubernetes
```bash
kubectl port-forward pod/podtest 8080:80
```
 
### 💻 Entrar en la línea de comandos del pod
```bash
kubectl exec -it podtest -- sh
```
 
### 📃 Ver logs de un pod
```bash
kubectl logs podtest -f
```
 
---
 
## 📄 Manifiestos de Kubernetes
 
Los YAML son la forma declarativa de definir objetos en Kubernetes.
 
📖 Plantilla oficial de pod: https://kubernetes.io/docs/concepts/workloads/pods/
 
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: podtest2
spec:
  containers:
  - name: cont1
    image: nginx:alpine
    # El template del pod termina aquí
```
 
### 🔎 Ver versiones y recursos de la API
```bash
kubectl api-versions
```
 
```bash
kubectl api-resources | grep Pod
```
 
---
 
## 🏷️ Labels
 
Dentro de `metadata` se asigna un array de labels.
 
> 💡 **Buena práctica:** asigna siempre al menos el label `app`.
 
```yaml
metadata:
  name: podtest2
  labels:
    app: front
    env: dev
```
 
### 🔍 Filtrar pods por label
```bash
kubectl get pods -l app=backend
```
 
---
 
## ⚠️ Problemas de los pods
 
- ❌ Sin **self-healing** (no se regeneran solos)
- ✍️ Crear pods masivamente requiere hacerlo manualmente en el YAML
- 🔄 Sin **auto-refresh**: no se actualizan solos
---
 
## 🧬 ReplicaSet
 
- 🔼 Objeto superior a los pods
- 🤝 Se "adueña" de ellos y los crea
- 🏷️ Agrega al `metadata` de los pods el valor `owner`, referenciando a qué ReplicaSet pertenecen
- 🚫 Otro ReplicaSet no puede tomar un pod que ya tenga owner
### 📋 Obtener ReplicaSets (shortname `rs`)
```bash
kubectl get rs
```
> ℹ️ Puedes consultar los shortnames disponibles con `kubectl api-resources`
 
### 🏷️ Agregar labels a pods sin owner
```bash
kubectl label pods podtest1 app=pod-label
```
 
> ⚠️ **Cuidado:** al crear dos pods diferentes con el mismo label, el ReplicaSet los adopta como suyos y agrega a `ownerReferences` la referencia del ReplicaSet — aunque los pods sean totalmente distintos entre sí, para el ReplicaSet son "iguales" por compartir label.
>
> Por eso, los pods deben crearse siempre mediante unidades u objetos superiores: **ReplicaSets** o **Deployments**.
 
### 🐛 Problemas de ReplicaSet
 
El ReplicaSet mantiene un número `n` de réplicas de un pod según lo definido en el manifiesto YAML.
 
Si se modifica un pod directamente (en caliente), **no ocurre nada**: el ReplicaSet solo vigila el número de pods que coinciden con el label definido en `metadata`, por lo que no puede cambiar los pods ni sus configuraciones.
 
---
 
## 🚢 Deployments
 
Un Deployment es un objeto que está por encima de un ReplicaSet, y este a su vez por encima del pod.
 
| Parámetro | Descripción | Valor por defecto |
|---|---|---|
| 🔽 `MaxAvailable` | Cuántos pods se permiten fuera de servicio | 25% |
| 🔼 `MaxSearch` | Pods adicionales permitidos al crear nuevos | — |
| 🗂️ Historial | ReplicaSets que Kubernetes mantiene por defecto | 10 |
 
### 🏷️ Mostrar labels de un deployment
```bash
kubectl get deployment --show-labels
```
 
### ✅ Verificar el éxito del rollout
```bash
kubectl rollout status deployment <nombreDeployment>
kubectl rollout status deployment deployment-test
```
 
### 🔗 OwnerReferences en Deployment
 
- Un **pod** tiene como `ownerReference` a un **ReplicaSet**
- Un **ReplicaSet** tiene como `ownerReference` a un **Deployment**
🚫 Esta jerarquía no puede saltarse: un pod nunca puede tener como `ownerReference` directamente a un Deployment.
 
---
 
## 🔄 Rolling Updates de Deployments
 
```bash
kubectl apply -f deployment.yaml
```
 
Al aplicar `apply -f` sobre el YAML del Deployment con algún cambio, este eliminará y creará pods con las nuevas especificaciones (según su estrategia configurada).
 
### ✅ Comprobar el estado del rollout
```bash
kubectl rollout status deployment <nombreDeployment>
```
 
### 🔬 Ver detalle con describe
```bash
kubectl describe deploy deployment-test
```
 
<details>
<summary>📋 Ejemplo de salida (Events)</summary>
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
</details>
---
 
## 📜 Histórico y revisiones de un Deployment
 
```bash
kubectl rollout history deployment deployment-test
```
Muestra las revisiones o rollouts ejecutados.
 
### 🏷️ Change-cause en un Deployment
 
Existen 3 maneras de definirlo:
 
**1️⃣ Flag `--record`** *(deprecated ⚠️, no se recomienda su uso)*
```bash
kubectl apply -f deployment.yaml --record
```
 
**2️⃣ Anotación en el YAML** ✅ *(recomendada)*
```yaml
metadata:
  annotations:
    # Anotación para definir una causa de cambio en el deployment
    kubernetes.io/change-cause: Changes port to 120
```
 
**3️⃣ Comando `annotate`**
```bash
kubectl annotate deployment.v1.apps/nginx-deployment kubernetes.io/change-cause="..."
```
 
### 🔎 Ver una revisión concreta
```bash
kubectl rollout history deploy deployment-test --revision=3
```
 
---
 
## ⏪ Rollbacks
 
Para volver a una versión anterior de un Deployment:
 
```bash
kubectl rollout undo deploy deployment-test --to-revision=3
```
 
> 💡 **Nota:** Kubernetes guarda por defecto hasta **10 revisiones** para poder volver atrás.
 
---
 
## 🌐 Servicios
 
Un servicio es un objeto que observa pods con cierto label (por ejemplo, `app=web`) y les proporciona:
 
- 🔒 Una **IP única** garantizada en el tiempo
- ⚖️ **Balanceo de carga** entre los pods disponibles (algoritmo de distribución aleatoria)
- 🌍 Un **DNS** consultable por el usuario
- 👀 Visibilidad sobre pods con cierto label, estén o no dentro de un ReplicaSet
### 🔗 Endpoints en un servicio
 
| | IP de Servicio | IP de Pod |
|---|---|---|
| Estabilidad | ✅ No cambia | ⚠️ Puede cambiar (el pod puede morir) |
 
El objeto **Endpoints** es una lista de IPs de los pods que cumplen el label del servicio:
 
- 🆕 Si nace un pod nuevo → se agrega su IP al endpoint
- ☠️ Si un pod muere → el servicio detecta la baja y elimina su IP del endpoint
De esta forma se mantiene la disponibilidad e integridad del tráfico.
 
### 🔬 Descripción de servicios
 
Por defecto, un servicio se crea de tipo `ClusterIP` (IP virtual):
 
```bash
kubectl describe svc my-service
```
 
<details>
<summary>📋 Ejemplo de salida</summary>
```
Name:              my-service
Namespace:         default
Labels:            app=front
Annotations:       <none>
Selector:          app=front
Type:              ClusterIP
IP Family Policy:  SingleStack
IP Families:       IPv4
IP:                10.101.207.125
IPs:               10.101.207.125
Port:              <unset>  8080/TCP
TargetPort:        80/TCP
Endpoints:         10.1.0.127:80,10.1.0.128:80,10.1.0.129:80
Session Affinity:  None
Events:            <none>
```
</details>
```bash
kubectl get po -l app=front -o wide
```
 
<details>
<summary>📋 Ejemplo de salida</summary>
```
NAME                               READY   STATUS    RESTARTS   AGE     IP           NODE             NOMINATED NODE   READINESS GATES
deployment-test-6cf85c55cf-mr254   1/1     Running   0          6m35s   10.1.0.129   docker-desktop   <none>           <none>
deployment-test-6cf85c55cf-s7xzc   1/1     Running   0          6m40s   10.1.0.127   docker-desktop   <none>           <none>
deployment-test-6cf85c55cf-wstdh   1/1     Running   0          6m38s   10.1.0.128   docker-desktop   <none>           <none>
```
</details>
> 🚫 **No es recomendable** crear pods fuera de ReplicaSets, como se ha comentado anteriormente.
 
### 🌍 Servicios y DNS
 
Cada servicio aporta su propio DNS. Se puede consultar por IP o por nombre DNS:
 
```bash
curl my-service:8080
curl <IP>:8080
```
 
### 🗂️ Tipos de servicios
 
| Tipo | Descripción |
|---|---|
| 🏠 **ClusterIP** | IP virtual asignada por Kubernetes, permanente en el tiempo. No accesible desde fuera del cluster. |
| 🚪 **NodePort** | Expone el servicio fuera del cluster a nivel de nodo. Rango de puertos por defecto: `30000-32767`. |
| ☁️ **LoadBalancer** | Kubernetes no proporciona balanceadores por defecto; se usan típicamente en entornos cloud. |
 
---
 
<div align="center">
📚 *Cheatsheet personal de Kubernetes — mantenido por Samuel* 
 
</div>