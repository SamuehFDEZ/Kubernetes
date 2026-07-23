<div align="center">

# ☸️ Kubernetes Cheatsheet

### Guía de referencia rápida: Pods, ReplicaSets, Deployments, Servicios, Namespaces y Recursos

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
12. [🐹 Hands-on: app Go en Docker + Kubernetes](#-hands-on-app-go-en-docker--kubernetes)
13. [🗂️ Namespaces](#️-namespaces)
14. [🧭 Contextos](#-contextos)
15. [📏 Límites de recursos (Limits & Requests)](#-límites-de-recursos-limits--requests)
16. [🎚️ Quality of Service (QoS)](#️-quality-of-service-qos)
17. [🚧 LimitRange](#-limitrange)
18. [📊 ResourceQuota](#-resourcequota)

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
Genera y muestra el manifiesto que se enviaría a la API, sin crear el objeto realmente. Es la forma más rápida de generar YAML de partida sin memorizar la sintaxis:
```bash
kubectl run nginx --image=nginx --dry-run=client -o yaml
```
> 💡 **Truco:** combina `--dry-run=client -o yaml` con `> pod.yaml` para generar plantillas base al vuelo y luego editarlas, en vez de escribir el YAML desde cero.

### ⚙️ Pod con spec sobrescrita mediante JSON
Permite sobrescribir campos puntuales del spec generado sin tener que reescribir todo el manifiesto:
```bash
kubectl run nginx --image=nginx --overrides='{ "apiVersion": "v1", "spec": { ... } }'
```

### 📦 Pod busybox interactivo, autolimpiable
`-i -t` (o `-it`) abren una sesión interactiva con TTY; `--rm` borra el pod en cuanto termina la sesión, muy útil para pods "de usar y tirar" de debug:
```bash
kubectl run -i -t busybox --image=busybox --restart=Never --rm
```

### 🎯 Comando por defecto con argumentos personalizados
Mantiene el `ENTRYPOINT` de la imagen pero le pasa argumentos distintos:
```bash
kubectl run nginx --image=nginx -- <arg1> <arg2> ... <argN>
```

### 🛠️ Comando y argumentos personalizados
Sustituye por completo el `ENTRYPOINT` de la imagen por el comando indicado:
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
Redirige un puerto local hacia el puerto de un pod, sin necesidad de exponer un Service. Ideal para depurar sin tocar la red del cluster:
```bash
kubectl port-forward pod/podtest 8080:80
```

### 💻 Entrar en la línea de comandos del pod
```bash
kubectl exec -it podtest -- sh
```

### 📃 Ver logs de un pod
El flag `-f` (`--follow`) mantiene el stream de logs abierto en tiempo real, como un `tail -f`:
```bash
kubectl logs podtest -f
```

---

## 📄 Manifiestos de Kubernetes

Los YAML son la forma **declarativa** de definir objetos en Kubernetes: describes el estado deseado y el control plane se encarga de converger hacia él, en lugar de ejecutar comandos imperativos paso a paso.

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

Dentro de `metadata` se asigna un array de labels: pares clave-valor que sirven para identificar y **seleccionar** objetos (por ReplicaSets, Services, `kubectl get -l`, etc.). No confundir con las *annotations*, que son solo metadatos informativos sin capacidad de selección.

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

- ❌ Sin **self-healing**: si un pod muere, nadie lo vuelve a crear automáticamente.
- ✍️ Crear pods masivamente requiere hacerlo manualmente, uno a uno, en el YAML.
- 🔄 Sin **auto-refresh**: si cambias la imagen o la config, el pod no se actualiza solo.

> Estos tres problemas son exactamente lo que resuelven los objetos de nivel superior (ReplicaSet y Deployment) que se explican a continuación.

---

## 🧬 ReplicaSet

- 🔼 Objeto superior a los pods: define cuántas réplicas de un pod deben existir en todo momento.
- 🤝 Se "adueña" de los pods que coinciden con su `selector` y los crea si faltan.
- 🏷️ Agrega a los pods una referencia (`ownerReferences`) indicando a qué ReplicaSet pertenecen.
- 🚫 Un pod que ya tiene owner no puede ser adoptado por otro ReplicaSet distinto.

### 📋 Obtener ReplicaSets (shortname `rs`)
```bash
kubectl get rs
```
> ℹ️ Puedes consultar todos los shortnames disponibles con `kubectl api-resources`.

### 🏷️ Agregar labels a pods sin owner
```bash
kubectl label pods podtest1 app=pod-label
```

> ⚠️ **Cuidado:** si creas dos pods distintos con el mismo label que usa el `selector` de un ReplicaSet, este los adopta como propios y les añade la referencia (`ownerReferences`) al ReplicaSet, aunque los pods sean completamente distintos entre sí — para el ReplicaSet son "iguales" por compartir label.
>
> Por eso, los pods deben crearse siempre a través de objetos superiores: **ReplicaSets** o **Deployments**, nunca sueltos.

### 🐛 Problemas de ReplicaSet

El ReplicaSet mantiene un número `n` de réplicas de un pod según lo definido en el manifiesto YAML.

Si se modifica un pod directamente (en caliente), **no ocurre nada visible en el ReplicaSet**: este solo vigila que el número de pods que coinciden con el `selector`/label definido en `metadata` sea el correcto, pero no puede alterar la configuración de los pods ya existentes. Para propagar cambios de configuración hace falta un Deployment.

---

## 🚢 Deployments

Un Deployment es un objeto que está por encima de un ReplicaSet, y este a su vez por encima del pod. Es el objeto que sí sabe gestionar actualizaciones de forma controlada (rolling updates, rollbacks, historial de revisiones).

| Parámetro | Descripción | Valor por defecto |
|---|---|---|
| 🔽 `maxUnavailable` | Cuántos pods se permiten fuera de servicio durante la actualización | 25% |
| 🔼 `maxSurge` | Pods adicionales por encima del número deseado permitidos al crear nuevos | 25% |
| 🗂️ Historial | ReplicaSets/revisiones que Kubernetes mantiene por defecto | 10 |

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

- Un **pod** tiene como `ownerReference` a un **ReplicaSet**.
- Un **ReplicaSet** tiene como `ownerReference` a un **Deployment**.

🚫 Esta jerarquía no puede saltarse: un pod nunca puede tener como `ownerReference` directamente a un Deployment.

---

## 🔄 Rolling Updates de Deployments

```bash
kubectl apply -f deployment.yaml
```

Al aplicar `apply -f` sobre el YAML del Deployment con algún cambio (por ejemplo, una nueva imagen), Kubernetes crea un **nuevo ReplicaSet** con la especificación actualizada y va escalándolo hacia arriba mientras escala el ReplicaSet antiguo hacia abajo, de forma gradual y controlada según `maxUnavailable`/`maxSurge` — no se destruyen todos los pods de golpe.

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

Un Service es un objeto que observa pods con cierto label (por ejemplo, `app=web`) y les proporciona:

- 🔒 Una **IP única** garantizada en el tiempo (aunque los pods por debajo cambien).
- ⚖️ **Balanceo de carga** entre los pods disponibles (algoritmo de distribución aleatoria por defecto).
- 🌍 Un **DNS** consultable por el usuario.
- 👀 Visibilidad sobre pods con cierto label, estén o no dentro de un ReplicaSet.

### 🔗 Endpoints en un servicio

| | IP de Servicio | IP de Pod |
|---|---|---|
| Estabilidad | ✅ No cambia | ⚠️ Puede cambiar (el pod puede morir y recrearse con otra IP) |

El objeto **Endpoints** es una lista de IPs de los pods que cumplen el `selector` del servicio:

- 🆕 Si nace un pod nuevo que cumple el label → se agrega su IP al endpoint.
- ☠️ Si un pod muere → el servicio detecta la baja y elimina su IP del endpoint.

De esta forma se mantiene la disponibilidad e integridad del tráfico sin que el cliente tenga que saber qué pods están vivos en cada momento.

### 🔬 Descripción de servicios

Por defecto, un servicio se crea de tipo `ClusterIP` (IP virtual, solo accesible dentro del cluster):

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

> 🚫 **No es recomendable** crear pods fuera de ReplicaSets/Deployments, como se ha comentado anteriormente.

### 🌍 Servicios y DNS

Cada servicio aporta su propio registro DNS interno (`<servicio>.<namespace>.svc.cluster.local`). Se puede consultar por IP o por nombre DNS:

```bash
curl my-service:8080
curl <IP>:8080
```

### 🗂️ Tipos de servicios

| Tipo | Descripción |
|---|---|
| 🏠 **ClusterIP** | IP virtual asignada por Kubernetes, permanente en el tiempo. No accesible desde fuera del cluster. Es el tipo por defecto. |
| 🚪 **NodePort** | Expone el servicio fuera del cluster a nivel de nodo, abriendo el mismo puerto en todos los nodos. Rango de puertos por defecto: `30000-32767`. |
| ☁️ **LoadBalancer** | Kubernetes no proporciona balanceadores por defecto; se implementa típicamente por el proveedor cloud (AWS ELB, GCP LB, etc.). En local (Docker Desktop) se simula asignando `localhost` como IP externa. |
| 🔗 **ExternalName** | No balancea pods: mapea el servicio a un nombre DNS externo mediante un `CNAME`. Útil para referenciar servicios fuera del cluster con un nombre interno consistente. |

---

## 🐹 Hands-on: app Go en Docker + Kubernetes

### 1️⃣ Contenedor de desarrollo para compilar/ejecutar la app en Go

**Linux:**
```bash
docker run --rm -dti -v $PWD/:/go --net host --name golang golang bash
```

| Flag | Significado |
|---|---|
| `docker run` | Crea y arranca un nuevo contenedor. |
| `--rm` | Al parar el contenedor, Docker lo elimina automáticamente (sin dejar contenedores parados como basura). |
| `-d` | Modo *detached*: el contenedor corre en segundo plano, sin quedar "enganchado" a la terminal. |
| `-t` | Asigna una pseudo-terminal (TTY) al contenedor. |
| `-i` | Modo interactivo: mantiene abierto el stdin aunque no estés conectado. |
| `-dti` | Es simplemente la combinación de esos tres flags juntos. |
| `-v $PWD/:/go` | Monta (*bind mount*) el directorio actual del host dentro del contenedor en `/go`. Todo lo que haya en la carpeta local es visible y editable desde el contenedor, y viceversa. `/go` es justo el `GOPATH` por defecto de la imagen oficial de Golang, por lo que es la forma típica de montar el código para compilarlo/ejecutarlo dentro. |
| `--net host` | El contenedor usa directamente la red del host en vez de tener red aislada (*bridge*). Si el proceso abre un puerto (p. ej. `:8080`), queda accesible directamente en `localhost:8080` sin mapear puertos con `-p`. ⚠️ Solo funciona así en **Linux**; en Docker Desktop para Windows/Mac tiene comportamiento limitado o distinto. |
| `--name golang` | Nombra el contenedor `golang`, para referenciarlo fácilmente (`docker exec -it golang ...`, `docker stop golang`, etc.) en lugar de usar el ID aleatorio. |
| `golang` | Imagen oficial de Go usada para crear el contenedor. |
| `bash` | Comando ejecutado al arrancar, en vez del entrypoint por defecto: abre una shell bash. |

> ⚠️ **Importante:** `$PWD` es sintaxis de Linux/macOS. En Windows cambia según la shell:

**CMD (Windows):**
```cmd
docker run --rm -it -v %cd%:/go -p 9090:9090 --name golang -w /go golang go run main.go
```
Esta variante arranca directamente la aplicación en Go (`-w /go` fija el *working directory* y ejecuta `go run main.go` de inmediato, sin abrir shell).

**PowerShell (Windows):**
```powershell
docker run --rm -dti -v ${PWD}:/go --net host --name golang golang bash
```

### 2️⃣ Crear el pod de prueba en Kubernetes

Forma correcta recomendada — con `--rm` el pod se autodestruye al salir de la sesión, evitando dejar pods sueltos de prueba:
```bash
kubectl run podtest3 --rm -ti --image=nginx:alpine -- sh
```

Si el pod ya existe y solo quieres entrar en él:
```bash
kubectl exec -it podtest3 -- sh
```

---

## 🗂️ Namespaces

Un Namespace es una **separación lógica** dentro del mismo cluster que permite organizar y aislar recursos.

- 📁 Ayuda a mantener orden separando recursos por proyecto o entorno (p. ej. un namespace de desarrollo y otro de pruebas de cliente).
- 🏢 Permite tener varios proyectos conviviendo en el mismo cluster sin pisarse.
- 🚧 Dentro de un namespace se pueden limitar: número de pods, hardware asignado (CPU, RAM, almacenamiento) y autorización/RBAC.

### 📋 Namespaces por defecto
```bash
kubectl get ns
# equivalente a:
kubectl get namespace
```

### ➕ Crear un namespace
```bash
kubectl create ns test-namespace
```

### 🌍 DNS de un servicio dentro de un namespace
El FQDN de un servicio sigue este patrón:
```
<nombre-servicio>.<nombre-namespace>.svc.cluster.local
```

---

## 🧭 Contextos

Un *contexto* combina cluster + usuario + namespace, y permite cambiar rápidamente el ámbito de trabajo de `kubectl` sin reescribir esos tres parámetros en cada comando.

### 🔎 Ver contextos disponibles
```bash
kubectl config view
```

### ➕ Crear un nuevo contexto
```bash
kubectl config set-context ci-context --namespace=ci --cluster=docker-desktop --user=Samuel
```

Al volver a listar la configuración veremos el nuevo contexto:
```bash
kubectl config view
```

### 🔄 Cambiar de contexto activo
```bash
kubectl config use-context ci-context
```
```
Switched to context "ci-context".
```

---

## 📏 Límites de recursos (Limits & Requests)

### Unidades

| Recurso | Unidades habituales |
|---|---|
| 💾 RAM | bytes, `Mi`/`Gi` (mebibytes/gibibytes) |
| ⚙️ CPU | `m` (milicores) — p. ej. `100m` = 0.1 CPU |

### Limits vs. Requests

- **`limits`**: techo máximo de recursos que el pod puede llegar a consumir. Por ejemplo, un límite de RAM de 30 Mi significa que el pod nunca podrá superar esos 30 Mi, por muy poco que use el nodo.
- **`requests`**: cantidad de recursos que se **reserva de forma garantizada** para el pod. Si un pod solicita `requests: 20Mi`, el scheduler solo lo colocará en un nodo que pueda garantizarle esos 20 Mi, y esa RAM queda apartada para él aunque no la use toda.

En otras palabras: el `request` es lo que el pod tiene garantizado, y el `limit` es hasta dónde se le permite crecer por encima de ese mínimo, siempre que el nodo tenga capacidad disponible. Si el pod se excede del `limit`, Kubernetes lo reinicia o lo elimina.

> 🔴 **`OOMKilled`**: estado que usa Kubernetes cuando un pod es terminado por quedarse sin memoria (*Out Of Memory*).
>
> 🟡 **`Pending`**: estado de un pod que está esperando a que algún nodo tenga capacidad suficiente de CPU/RAM para poder programarlo (*schedular*lo).

### ⚙️ Limitar recursos de CPU

Ver fichero de ejemplo `limit-cpu.yaml`.

Para diagnosticar los recursos que está usando un nodo del cluster:
```bash
kubectl describe nodes <nombre-del-nodo>
```

---

## 🎚️ Quality of Service (QoS)

Kubernetes asigna automáticamente una clase de QoS a cada pod según cómo defina sus `requests`/`limits`. Esta clase determina qué pods se eliminan primero cuando el nodo se queda sin recursos.

| Clase | Condición | Descripción |
|---|---|---|
| 🟢 **Guaranteed** | `requests` == `limits` en CPU y RAM, en todos los contenedores | Máxima prioridad; los últimos en ser eliminados ante falta de recursos. |
| 🟡 **Burstable** | Tiene `requests`/`limits` definidos, pero no son iguales | Trabajan normalmente por debajo de lo solicitado, pudiendo exceder ese uso puntualmente hasta el límite definido. |
| 🔴 **BestEffort** | No define ningún `request` ni `limit` | Los más peligrosos: sin ningún tope, pueden consumir recursos hasta agotar un nodo. Son los primeros en ser eliminados si hace falta liberar recursos. |

---

## 🚧 LimitRange

El `LimitRange` actúa **a nivel de objeto individual**, dentro de un namespace: valida que cada pod/contenedor que se cree respete los mínimos y máximos de recursos definidos en su manifiesto YAML.

- ✅ Solo aplica a los objetos que estén dentro del namespace donde se definió el `LimitRange`. Un pod sin namespace indicado no tendrá `LimitRange` asociado.
- 🚫 Si se intenta crear un pod que **supera el máximo** permitido, la API lo rechaza:
```
Error from server (Forbidden): error when creating "minMaxLimits.yaml": pods "podtest4" is forbidden: [maximum cpu usage per Container is 1, but limit is 2, maximum memory usage per Container is 1Gi, but limit is 2G]
```
- 🚫 Si el pod **no llega al mínimo** exigido de CPU o RAM, también es rechazado:
```
The Pod "podtest4" is invalid:
* spec.containers[0].resources.requests: Invalid value: "300m": must be less than or equal to cpu limit of 50m
* spec.containers[0].resources.requests: Invalid value: "400M": must be less than or equal to memory limit of 50M
```

---

## 📊 ResourceQuota

El `ResourceQuota` aplica **a nivel de namespace completo**, a diferencia del `LimitRange` que valida objeto por objeto.

- 📐 Limita la **suma total** de todos los recursos individuales consumidos dentro del namespace (no entiende de objetos concretos, solo del agregado).
- 🤝 No sustituye al `LimitRange`: son complementarios. El `LimitRange` valida que cada pod individual sea razonable; el `ResourceQuota` valida que la suma de todos ellos no desborde el presupuesto del namespace.

En un resourceQuota hemos de definir tanto el request como el limit para garantizar que no se pasen los pods de lo establecido en el manifiesto

Además de limitar el uso de ram y cpu también podemos controlar el numero de objetos que queremos tener en un namespace, por ejemplo, un namespace será capaz únicamente de crear pods

---

## Probe

Es un estado que se ejecuta para comprobar que el contenedor se encuentra en buen estado

- Probe es un diagnostico realizado por kubelet el cual corre en cada nodo, es el encargado de realizar los diagnosticos sobre los contenedores, dado un contenedor le asignamos un probe y un rango de tiempo, dado esto, kubelet irá preguntando al contenedor si está correcto, en caso contrario, tomará una acción contra el contenedor

Como pregunta kubelet
    - Mediante comando: kubelet ejecuta un comando, si devuelve 0 está ok, en caso contrario, ko
    - Por TCP: si el puerto funciona y responde todo ok, si el puerto no responde kubelet asume que hay un error
    - Http: kubelet hace una peticion yaml hacia el contenedor, cualquier codigo de 200 es ok, entre 400 y 500 algo hay mal

Tipos de probe
 - Liveness: si la aplicacion esta viva
   - Es una prueba que kubelet ejecuta en el contenedor cada x segundos, en esta prueba solo esperamos una resuesta del contenedor, si por ejemplo tenemos una web, basta con una peticion GET, si que puede ocurrir que para el contenedor la aplicacion está ok pero nosotros si accedemos puede haber un 500, por ello, con un liveness nos aseguramos de reinicar la aplicacion si no está funcionando bien
 - readiness: para ver si la aplicacion esta lista
   - Tenemos un servicio con dos pods, digamos que queremos agregar un nuevo pod pero queremos garantizar que cuando el pod esté listo pueda ponerse a servicio, para eso sirve el readiness, una especie de endpoint validador para saber si un pod está listo o no para recibir peticiones, si no pasa el readiness no se incluye
 -  Startup: para aplicaciones que tardan en arrancar o inciarse
   - El startup se usa para aplicaciones grandes, validando que el pod, hasta que no esté todo desplegado y todo listo no se despliegue

## ConfigMaps y variables de entorno

Para las variables de entorno las crearemos de la siguiente forma, para mas detalle consultar yaml en /envs/env.yaml

      env:
        - name: VAR1
          value: "valor de prueba 2"
        - name: VAR2
          value: "TEST2"
        - name: VAR3
          value: "TEST3"

Pero basicamente es esto, luego esas variables quedarán referenciadas dentro del pod, si mediante bash nos introducimos en el y escribimos env dentro del pod para ver las variables que almacena

Además otra forma, vista en el fichero /env/ref.yaml es que podemos hacer referencia a las variables en base a los manifiestos que k8s crea de sus objetos, por ejemplo:

      env:
        - name: MY_NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName

Ese spec.NodeName lo podemos observar mediante un describe -o yaml del pod, es una forma de recuperar las variables de manera dinámica como hariamos en cualquier programa

### ConfigMap

Es otro objeto de Kubernetes que se puede crear actualizar y eliminar, un configMap se creó para separar las configuraciones y hacer mas portables las configurciones de los pods, se usa también para evitar hardcodear en los pods, el pod consumirá todo lo que el configMap tenga definido como la imagen de un nginx

El configMap se vasa en Key Value, en el pod se referenciará a la Key del configMap, así consumirá a la vez los valores de esa Key dentro del configMap

Una forma de crear un configMap es, mediante un archivo .conf pasarlo a un configMap, la forma más idonea es de hacerlo en un manifiesto de kubernetes donde escribiremos las configuraciones de un configMap

El pod podrá ver el configMap mediante variables de entorno que hemos visto antes

Por defecto un pod toma la configuracion del servicio que aloje, en este caso, el de nginx

kubectl create configmap nginx-config --from-file=.\configMapExamples\nginx.conf

kubectl get cm

Dado el archivo cm-nginx.yaml si eliminamos las keys del configmap 

## Secret

Nos ayuda a guardar informacion sensible que no deberia ser visible para todo el mundo, ya sean tokens o contraseñas, a diferencia de los configMap donde podemos guardar por ejemplo nombres de configuracion

El secret está aislado del pod, si se modifica el secret, el cambio se aplica en el pod, siendo más facil de configurar

Desde un pod se accede a un secret como un configMap, mediante variables env o mediante volume

Para crear un secret

kubectl create secret generic mitesoro --from-file=.\secret-files\test.txt

kubectl describe secrets mitesoro

cuando hacemos un describe kubernetes no nos muestra el valor, pero si lo obtenemos mediante 
kubectl get secrets mitesoro -o yaml

tenemos los datos en base64


Una herramienta para reemplazar secrets para fortalezer la seguridad es envsubst

ensubst < secure.yaml > tmp.yaml

hacemos un cat tmp.yaml

y ahí es donde tenemos los datos reales

con secure.yaml lo tenemos encodeado

<div align="center">
📚 *Cheatsheet personal de Kubernetes — mantenido por Samuel*

</div>