<div align="center">
# ⚡ kubectl Mega Cheatsheet
 
### Todos los comandos de `kubectl`, atajos, trucos + plantilla completa de objetos Kubernetes
 
---
 
📦 `kubectl` · 🎯 Shortnames · ⚙️ `kubeconfig` · 📄 Manifiestos YAML de todos los objetos
 
</div>
---
 
## 📑 Índice
 
1. [⚙️ Configuración inicial y autocompletado](#️-configuración-inicial-y-autocompletado)
2. [🧭 Contextos, clusters y kubeconfig](#-contextos-clusters-y-kubeconfig)
3. [🗂️ Namespaces](#️-namespaces-1)
4. [🔍 kubectl get — consulta de recursos](#-kubectl-get--consulta-de-recursos)
5. [🔬 describe, logs y depuración](#-describe-logs-y-depuración)
6. [✏️ Crear, aplicar y editar](#️-crear-aplicar-y-editar)
7. [🗑️ Borrar recursos](#️-borrar-recursos)
8. [💻 exec, cp, port-forward y proxy](#-exec-cp-port-forward-y-proxy)
9. [🏷️ Labels y annotations](#️-labels-y-annotations)
10. [📈 Escalado y autoescalado](#-escalado-y-autoescalado)
11. [🚢 Rollouts (Deployments)](#-rollouts-deployments)
12. [🔐 RBAC y ServiceAccounts](#-rbac-y-serviceaccounts)
13. [📊 Recursos, nodos y cluster](#-recursos-nodos-y-cluster)
14. [🎩 Atajos y trucos avanzados](#-atajos-y-trucos-avanzados)
15. [📇 Tabla de shortnames](#-tabla-de-shortnames)
16. [🧩 MEGA PLANTILLA: todos los objetos de Kubernetes](#-mega-plantilla-todos-los-objetos-de-kubernetes)
---
 
## ⚙️ Configuración inicial y autocompletado
 
### Alias `k` (el truco más usado por cualquier admin de k8s)
```bash
alias k=kubectl
complete -o default -F __start_kubectl k
```
> 💡 Añádelo a tu `.bashrc`/`.zshrc`. Te ahorra miles de pulsaciones a la semana.
 
### Autocompletado de kubectl
```bash
# Bash
source <(kubectl completion bash)
echo 'source <(kubectl completion bash)' >> ~/.bashrc
 
# Zsh
source <(kubectl completion zsh)
echo 'source <(kubectl completion zsh)' >> ~/.zshrc
 
# PowerShell
kubectl completion powershell | Out-String | Invoke-Expression
```
 
### Ver versión de cliente y servidor
```bash
kubectl version
kubectl version --client
```
 
### Ayuda de cualquier comando/recurso
```bash
kubectl help
kubectl explain pod
kubectl explain pod.spec.containers      # 🎯 Muy útil: documentación de campos anidados
kubectl explain deployment.spec.strategy --recursive
```
 
---
 
## 🧭 Contextos, clusters y kubeconfig
 
```bash
kubectl config view                                  # Ver kubeconfig completo
kubectl config view --minify                         # Solo el contexto activo
kubectl config get-contexts                           # Listar contextos
kubectl config current-context                        # Contexto activo
kubectl config use-context <nombre>                    # Cambiar de contexto
kubectl config set-context <nombre> --namespace=<ns> --cluster=<cluster> --user=<user>
kubectl config set-context --current --namespace=<ns>  # 🎯 Cambia el namespace por defecto del contexto actual
kubectl config delete-context <nombre>
kubectl cluster-info                                   # Endpoints del cluster
kubectl cluster-info dump                              # Volcado completo de estado (debug profundo)
```
 
> 💡 **Truco:** en vez de escribir `-n <namespace>` en cada comando, usa `kubectl config set-context --current --namespace=<ns>` una vez y todos los comandos posteriores lo usarán por defecto.
 
---
 
## 🗂️ Namespaces
 
```bash
kubectl get ns
kubectl create ns <nombre>
kubectl delete ns <nombre>
kubectl get all -n <namespace>              # Todo lo que hay en un namespace
kubectl get all --all-namespaces            # Todo, en todos los namespaces
kubectl get all -A                          # 🎯 -A es el atajo de --all-namespaces
```
 
---
 
## 🔍 kubectl get — consulta de recursos
 
```bash
kubectl get pods                            # Pods del namespace actual
kubectl get pods -o wide                    # + IP, nodo, etc.
kubectl get pods -A                         # De todos los namespaces
kubectl get pods --show-labels              # Ver labels de cada pod
kubectl get pods -l app=web,env=prod        # Filtrar por varios labels (AND)
kubectl get pods -l 'env in (prod,staging)' # Filtro con selector de conjunto
kubectl get pods --field-selector=status.phase=Running
kubectl get pods --sort-by=.metadata.creationTimestamp   # 🎯 Ordenar por antigüedad
kubectl get pods --sort-by=.status.containerStatuses[0].restartCount
kubectl get pods -o yaml
kubectl get pods -o json
kubectl get pods -o name                    # Solo nombres (útil para scripting)
kubectl get pod <nombre> -o jsonpath='{.status.podIP}'   # 🎯 Extraer un campo concreto
kubectl get pods --watch                    # Ver cambios en tiempo real
kubectl get pods -w                         # Atajo de --watch
kubectl get events --sort-by='.lastTimestamp'             # Eventos del cluster, cronológicos
kubectl get events -A --field-selector type=Warning       # 🎯 Solo warnings, todos los ns
```
 
> 💡 **Truco `-o custom-columns`:** para tablas a medida sin depender de `-o wide`:
> ```bash
> kubectl get pods -o custom-columns='NAME:.metadata.name,IP:.status.podIP,NODE:.spec.nodeName'
> ```
 
---
 
## 🔬 describe, logs y depuración
 
```bash
kubectl describe pod <nombre>               # Detalle + eventos del pod
kubectl describe node <nombre>              # Recursos y condiciones del nodo
kubectl logs <pod>                          # Logs del pod
kubectl logs <pod> -f                       # Logs en streaming (follow)
kubectl logs <pod> --previous               # 🎯 Logs del contenedor anterior (tras un crash)
kubectl logs <pod> -c <contenedor>          # Logs de un contenedor concreto (pod multi-container)
kubectl logs -l app=web --all-containers    # Logs de todos los pods que cumplan el label
kubectl logs -l app=web --prefix            # 🎯 Antepone el nombre del pod a cada línea
kubectl top pod                             # Consumo de CPU/RAM por pod (requiere metrics-server)
kubectl top node                            # Consumo de CPU/RAM por nodo
kubectl get pod <nombre> -o yaml | less     # Inspección rápida del spec aplicado
kubectl debug <pod> -it --image=busybox     # 🎯 Contenedor de debug efímero adjunto al pod
```
 
---
 
## ✏️ Crear, aplicar y editar
 
```bash
kubectl apply -f manifiesto.yaml            # ✅ Forma declarativa recomendada
kubectl apply -f ./carpeta/                 # Aplica todos los YAML de una carpeta
kubectl apply -f ./carpeta/ -R              # Recursivo en subcarpetas
kubectl apply -f https://raw.githubusercontent.com/.../manifest.yaml   # Directo desde URL
kubectl create -f manifiesto.yaml           # Forma imperativa (falla si ya existe)
kubectl create deployment web --image=nginx --replicas=3
kubectl create configmap cfg --from-literal=key=value
kubectl create secret generic sec --from-literal=password=1234
kubectl edit deployment <nombre>            # Editor interactivo en vivo
kubectl patch deployment <nombre> -p '{"spec":{"replicas":5}}'   # 🎯 Cambios puntuales sin editor
kubectl replace -f manifiesto.yaml          # Sustituye el objeto entero
kubectl diff -f manifiesto.yaml             # 🎯 Muestra qué cambiaría un apply, sin aplicarlo
kubectl apply -f manifiesto.yaml --dry-run=server   # Valida contra la API sin persistir
```
 
> 💡 **Truco:** genera plantillas base sin memorizar YAML:
> ```bash
> kubectl create deployment web --image=nginx --dry-run=client -o yaml > deployment.yaml
> ```
 
---
 
## 🗑️ Borrar recursos
 
```bash
kubectl delete pod <nombre>
kubectl delete -f manifiesto.yaml
kubectl delete pods --all                   # Todos los pods del namespace
kubectl delete pods -l app=web              # Por label
kubectl delete pod <nombre> --grace-period=0 --force   # ⚠️ Forzar borrado inmediato (usar con cuidado)
kubectl delete ns <nombre>                  # Borra el namespace y TODO su contenido
```
 
---
 
## 💻 exec, cp, port-forward y proxy
 
```bash
kubectl exec -it <pod> -- sh                          # Shell interactiva
kubectl exec -it <pod> -c <contenedor> -- bash         # En un contenedor concreto
kubectl exec <pod> -- env                              # Ejecutar comando puntual sin sesión
kubectl cp <pod>:/ruta/archivo ./local                 # Copiar archivo del pod al host
kubectl cp ./local <pod>:/ruta/archivo                 # Copiar del host al pod
kubectl port-forward pod/<nombre> 8080:80
kubectl port-forward svc/<nombre> 8080:80              # 🎯 También funciona sobre Services
kubectl port-forward deployment/<nombre> 8080:80       # Y sobre Deployments
kubectl proxy --port=8001                              # Proxy local hacia la API server
```
 
---
 
## 🏷️ Labels y annotations
 
```bash
kubectl label pod <nombre> app=web                     # Añadir label
kubectl label pod <nombre> app=web --overwrite          # Sobrescribir uno existente
kubectl label pod <nombre> app-                          # 🎯 El "-" al final ELIMINA el label
kubectl annotate pod <nombre> descripcion="texto"
kubectl annotate pod <nombre> descripcion-                # Eliminar annotation
```
 
---
 
## 📈 Escalado y autoescalado
 
```bash
kubectl scale deployment <nombre> --replicas=5
kubectl scale rs <nombre> --replicas=3
kubectl autoscale deployment <nombre> --min=2 --max=10 --cpu-percent=80   # Crea un HPA
kubectl get hpa
```
 
---
 
## 🚢 Rollouts (Deployments)
 
```bash
kubectl rollout status deployment <nombre>
kubectl rollout history deployment <nombre>
kubectl rollout history deployment <nombre> --revision=3
kubectl rollout undo deployment <nombre>                     # Vuelve a la revisión anterior
kubectl rollout undo deployment <nombre> --to-revision=3      # Vuelve a una revisión concreta
kubectl rollout restart deployment <nombre>                   # 🎯 Reinicia todos los pods sin cambiar el YAML (útil tras rotar un secret)
kubectl rollout pause deployment <nombre>                     # Pausa el rollout (para agrupar varios cambios)
kubectl rollout resume deployment <nombre>
```
 
---
 
## 🔐 RBAC y ServiceAccounts
 
```bash
kubectl create sa <nombre>
kubectl create role <nombre> --verb=get,list,watch --resource=pods
kubectl create rolebinding <nombre> --role=<role> --serviceaccount=<ns>:<sa>
kubectl create clusterrole <nombre> --verb=get,list --resource=pods
kubectl create clusterrolebinding <nombre> --clusterrole=<cr> --serviceaccount=<ns>:<sa>
kubectl auth can-i create pods                          # 🎯 ¿Puedo yo crear pods?
kubectl auth can-i create pods --as=system:serviceaccount:<ns>:<sa>   # ¿Puede esa SA?
kubectl auth can-i '*' '*'                              # ¿Soy cluster-admin?
```
 
---
 
## 📊 Recursos, nodos y cluster
 
```bash
kubectl get nodes
kubectl describe node <nombre>
kubectl cordon <nodo>                        # Marca el nodo como no-programable
kubectl drain <nodo> --ignore-daemonsets       # 🎯 Vacía un nodo antes de mantenimiento
kubectl uncordon <nodo>                        # Vuelve a permitir programar en él
kubectl taint nodes <nodo> key=value:NoSchedule   # Añadir taint
kubectl taint nodes <nodo> key:NoSchedule-         # Eliminar taint
kubectl api-resources                           # Todos los tipos de recursos disponibles
kubectl api-resources --namespaced=true         # Solo los que van por namespace
kubectl api-versions
```
 
---
 
## 🎩 Atajos y trucos avanzados
 
| Truco | Comando |
|---|---|
| 🎯 Salida solo con nombres, para pipes con `xargs` | `kubectl get pods -o name` |
| 🎯 Borrar todos los pods en estado `Evicted` | `kubectl get pods \| grep Evicted \| awk '{print $1}' \| xargs kubectl delete pod` |
| 🎯 Reiniciar un pod sin borrar el Deployment | `kubectl delete pod <nombre>` *(el ReplicaSet lo recrea solo)* |
| 🎯 Ver el YAML "real" aplicado, sin campos generados | `kubectl get pod <nombre> -o yaml --export` *(deprecado; usar `-o yaml` y limpiar manualmente en versiones nuevas)* |
| 🎯 Simular un `apply` para ver el diff antes de aplicar | `kubectl diff -f manifiesto.yaml` |
| 🎯 Ejecutar un pod temporal para hacer un `curl` de prueba | `kubectl run tmp --rm -it --image=curlimages/curl --restart=Never -- sh` |
| 🎯 Ver todos los eventos de tipo `Warning` del cluster | `kubectl get events -A --field-selector type=Warning` |
| 🎯 Forzar el borrado de un namespace atascado en `Terminating` | editar el `finalizers: []` vía `kubectl proxy` + `curl` al API |
| 🎯 Ver qué pod consume más memoria | `kubectl top pod --sort-by=memory` |
| 🎯 Ver qué pod consume más CPU | `kubectl top pod --sort-by=cpu` |
| 🎯 Contexto + namespace en el prompt de la terminal | usar `kubectx` / `kubens` (herramientas externas muy populares) |
| 🎯 Ver el manifiesto sin metadata "ruidosa" (uid, resourceVersion...) | `kubectl neat` (plugin de `krew`) |
| 🎯 Instalar plugins de kubectl | `kubectl krew install <plugin>` (requiere instalar `krew` antes) |
 
---
 
## 📇 Tabla de shortnames
 
| Recurso completo | Shortname |
|---|---|
| `namespaces` | `ns` |
| `nodes` | `no` |
| `pods` | `po` |
| `replicationcontrollers` | `rc` |
| `replicasets` | `rs` |
| `deployments` | `deploy` |
| `statefulsets` | `sts` |
| `daemonsets` | `ds` |
| `services` | `svc` |
| `endpoints` | `ep` |
| `configmaps` | `cm` |
| `persistentvolumes` | `pv` |
| `persistentvolumeclaims` | `pvc` |
| `storageclasses` | `sc` |
| `horizontalpodautoscalers` | `hpa` |
| `ingresses` | `ing` |
| `networkpolicies` | `netpol` |
| `serviceaccounts` | `sa` |
| `clusterroles` | — *(sin shortname)* |
| `clusterrolebindings` | — *(sin shortname)* |
| `customresourcedefinitions` | `crd`, `crds` |
| `events` | `ev` |
| `limitranges` | `limits` |
| `resourcequotas` | `quota` |
| `jobs` | — *(sin shortname)* |
| `cronjobs` | `cj` |
 
> ℹ️ Lista completa siempre actualizada con: `kubectl api-resources`
 
---
 
## 🧩 MEGA PLANTILLA: todos los objetos de Kubernetes
 
> 📌 Cada bloque es un manifiesto independiente y funcional (con valores de ejemplo). Puedes copiar solo el que necesites o guardarlos todos en un mismo archivo separados por `---` y aplicarlos juntos con `kubectl apply -f mega-plantilla.yaml`.
 
### 🧱 Pod
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: pod-ejemplo
  namespace: default
  labels:
    app: ejemplo
spec:
  containers:
    - name: app
      image: nginx:alpine
      ports:
        - containerPort: 80
      env:
        - name: ENV_VAR
          value: "valor"
      resources:
        requests:
          cpu: "100m"
          memory: "128Mi"
        limits:
          cpu: "250m"
          memory: "256Mi"
```
 
### 🧬 ReplicaSet
```yaml
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: rs-ejemplo
spec:
  replicas: 3
  selector:
    matchLabels:
      app: ejemplo
  template:
    metadata:
      labels:
        app: ejemplo
    spec:
      containers:
        - name: app
          image: nginx:alpine
```
 
### 🚢 Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: deploy-ejemplo
  annotations:
    kubernetes.io/change-cause: "Version inicial"
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 25%
      maxSurge: 25%
  selector:
    matchLabels:
      app: ejemplo
  template:
    metadata:
      labels:
        app: ejemplo
    spec:
      containers:
        - name: app
          image: nginx:alpine
          ports:
            - containerPort: 80
```
 
### 🗄️ StatefulSet
Para cargas con estado (bases de datos, colas): identidad de red estable y almacenamiento persistente por réplica.
```yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: sts-ejemplo
spec:
  serviceName: "sts-ejemplo-headless"
  replicas: 3
  selector:
    matchLabels:
      app: sts-ejemplo
  template:
    metadata:
      labels:
        app: sts-ejemplo
    spec:
      containers:
        - name: app
          image: postgres:16
          ports:
            - containerPort: 5432
          volumeMounts:
            - name: data
              mountPath: /var/lib/postgresql/data
  volumeClaimTemplates:
    - metadata:
        name: data
      spec:
        accessModes: ["ReadWriteOnce"]
        resources:
          requests:
            storage: 1Gi
```
 
### 🛰️ DaemonSet
Garantiza un pod por cada nodo del cluster (agentes de logging, monitorización, red).
```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: ds-ejemplo
spec:
  selector:
    matchLabels:
      app: agente
  template:
    metadata:
      labels:
        app: agente
    spec:
      containers:
        - name: agente
          image: fluent/fluent-bit
```
 
### 🕐 Job
Ejecuta un pod hasta completarlo (tarea puntual, no un servicio continuo).
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: job-ejemplo
spec:
  backoffLimit: 3
  template:
    spec:
      restartPolicy: Never
      containers:
        - name: tarea
          image: busybox
          command: ["sh", "-c", "echo Procesando... && sleep 5"]
```
 
### ⏰ CronJob
Igual que un Job, pero programado periódicamente con sintaxis cron.
```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: cronjob-ejemplo
spec:
  schedule: "*/5 * * * *"       # Cada 5 minutos
  jobTemplate:
    spec:
      template:
        spec:
          restartPolicy: OnFailure
          containers:
            - name: tarea
              image: busybox
              command: ["sh", "-c", "echo Ejecutando tarea programada"]
```
 
### 🌐 Service — ClusterIP (por defecto)
```yaml
apiVersion: v1
kind: Service
metadata:
  name: svc-clusterip
spec:
  type: ClusterIP
  selector:
    app: ejemplo
  ports:
    - port: 8080
      targetPort: 80
```
 
### 🚪 Service — NodePort
```yaml
apiVersion: v1
kind: Service
metadata:
  name: svc-nodeport
spec:
  type: NodePort
  selector:
    app: ejemplo
  ports:
    - port: 8080
      targetPort: 80
      nodePort: 30080     # Rango válido: 30000-32767
```
 
### ☁️ Service — LoadBalancer
```yaml
apiVersion: v1
kind: Service
metadata:
  name: svc-loadbalancer
spec:
  type: LoadBalancer
  selector:
    app: ejemplo
  ports:
    - port: 80
      targetPort: 80
```
 
### 🔗 Service — ExternalName
```yaml
apiVersion: v1
kind: Service
metadata:
  name: svc-externalname
spec:
  type: ExternalName
  externalName: api.miservicio-externo.com
```
 
### 🌍 Ingress
Enrutado HTTP/HTTPS hacia distintos Services según host/path (requiere un Ingress Controller instalado en el cluster).
```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress-ejemplo
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  ingressClassName: nginx
  rules:
    - host: app.ejemplo.local
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: svc-clusterip
                port:
                  number: 8080
  tls:
    - hosts:
        - app.ejemplo.local
      secretName: tls-ejemplo
```
 
### 🗒️ ConfigMap
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm-ejemplo
data:
  APP_ENV: "production"
  config.properties: |
    timeout=30
    retries=3
```
 
### 🔑 Secret
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: secret-ejemplo
type: Opaque
stringData:                 # stringData acepta texto plano (k8s lo codifica en base64 solo)
  DB_PASSWORD: "SuperSecreta123"
  API_KEY: "abc123xyz"
```
 
### 💾 PersistentVolume (PV)
```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-ejemplo
spec:
  capacity:
    storage: 5Gi
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Retain
  storageClassName: manual
  hostPath:
    path: /mnt/data
```
 
### 📌 PersistentVolumeClaim (PVC)
```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pvc-ejemplo
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 5Gi
  storageClassName: manual
```
 
### 🏗️ StorageClass
Define cómo se aprovisiona el almacenamiento dinámicamente (normalmente lo aporta el proveedor cloud).
```yaml
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: sc-ejemplo
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp3
reclaimPolicy: Delete
volumeBindingMode: WaitForFirstConsumer
```
 
### 🗂️ Namespace
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: ns-ejemplo
  labels:
    entorno: desarrollo
```
 
### 🚧 LimitRange
```yaml
apiVersion: v1
kind: LimitRange
metadata:
  name: limitrange-ejemplo
  namespace: ns-ejemplo
spec:
  limits:
    - type: Container
      default:
        cpu: "500m"
        memory: "512Mi"
      defaultRequest:
        cpu: "250m"
        memory: "256Mi"
      max:
        cpu: "1"
        memory: "1Gi"
      min:
        cpu: "50m"
        memory: "64Mi"
```
 
### 📊 ResourceQuota
```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: quota-ejemplo
  namespace: ns-ejemplo
spec:
  hard:
    requests.cpu: "4"
    requests.memory: 4Gi
    limits.cpu: "8"
    limits.memory: 8Gi
    pods: "20"
```
 
### 📈 HorizontalPodAutoscaler (HPA)
```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: hpa-ejemplo
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: deploy-ejemplo
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 75
```
 
### 🛡️ NetworkPolicy
Controla qué tráfico entra/sale de los pods seleccionados (requiere un CNI compatible, como Calico).
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: netpol-ejemplo
spec:
  podSelector:
    matchLabels:
      app: ejemplo
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: frontend
      ports:
        - protocol: TCP
          port: 80
  egress:
    - to:
        - podSelector:
            matchLabels:
              app: base-datos
      ports:
        - protocol: TCP
          port: 5432
```
 
### 👤 ServiceAccount
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: sa-ejemplo
  namespace: ns-ejemplo
```
 
### 🔐 Role + RoleBinding (permisos dentro de un namespace)
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: role-ejemplo
  namespace: ns-ejemplo
rules:
  - apiGroups: [""]
    resources: ["pods", "pods/log"]
    verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: rolebinding-ejemplo
  namespace: ns-ejemplo
subjects:
  - kind: ServiceAccount
    name: sa-ejemplo
    namespace: ns-ejemplo
roleRef:
  kind: Role
  name: role-ejemplo
  apiGroup: rbac.authorization.k8s.io
```
 
### 🔐 ClusterRole + ClusterRoleBinding (permisos a nivel de cluster)
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: clusterrole-ejemplo
rules:
  - apiGroups: [""]
    resources: ["nodes"]
    verbs: ["get", "list", "watch"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: clusterrolebinding-ejemplo
subjects:
  - kind: ServiceAccount
    name: sa-ejemplo
    namespace: ns-ejemplo
roleRef:
  kind: ClusterRole
  name: clusterrole-ejemplo
  apiGroup: rbac.authorization.k8s.io
```
 
### 🎯 PodDisruptionBudget (PDB)
Garantiza un mínimo de pods disponibles durante mantenimientos voluntarios (`drain`, actualizaciones de nodo).
```yaml
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: pdb-ejemplo
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app: ejemplo
```
 
### 🧩 CustomResourceDefinition (CRD)
Extiende la API de Kubernetes con tipos de recurso propios (base de los operadores).
```yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: widgets.ejemplo.com
spec:
  group: ejemplo.com
  versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                tamano:
                  type: string
  scope: Namespaced
  names:
    plural: widgets
    singular: widget
    kind: Widget
    shortNames: ["wg"]
```
 
---
 
<div align="center">
📚 *kubectl Mega Cheatsheet — mantenido por Samuel*
</div>