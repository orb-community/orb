# Orb local k8s cluster

Follow those steps to setup a local k8s cluster and deploy Orb.

## 🧱 Requirements

- [Docker Environment](#docker)
- [Helm 3](#helm-3)
- [Kubectl](#kubectl)
- [Kind](#install-kind)

> **💡 Note:** If you have those installed, please skip to [Deploy Orb on Kind](#deploy-orb-kind).

<a name="docker"></a>
## 🐳 Docker Environment (Requirement)

Quick install a **Docker** executing:
```shell
make install-docker
```

Check if you have a **Docker** running by executing:
```shell
docker version
```
If you need help to setup a **Docker Environment**, follow the [steps from here](https://docs.docker.com/engine/install/debian/).

<a name="helm-3"></a>
## ⚓ Helm 3 (Requirement)

Quick install a **Helm 3** executing:
```shell
make install-helm
```

Check if you have **Helm 3** installed by executing:
```shell
helm version
```
If you need help to install **Helm 3**, follow the [steps from here](https://helm.sh/docs/intro/install/).

> 🚨 **Warning:** Make sure you have version 3 installed, orb helm charts doesn't officialy support helm 2.

<a name="kubectl"></a>
## 🐋 Kubectl (Requirement)

Quick install a **Kubectl** executing:
```shell
make install-kubectl
```

Check if you have **Kubectl** cmd installed by executing:
```shell
kubectl version --client
```
If you need help to install **Kubectl**, follow the [steps from here](https://kubernetes.io/docs/tasks/tools/).

<a name="install-kind"></a>
## 🚢 Install Kind (Requirement)

Kind is a tool for running local k8s clusters using docker container as nodes.

Quick install a **Kind** on Linux executing:
```shell
make install-kind
```

If you have `go 1.17 or later` installed:
```shell
go install sigs.k8s.io/kind@v0.14.0
```

macOS users can also use `brew`:
```shell
brew install kind
```

> 🚨 **Windows WSL users**: WSL is also supported, but for some reason the Orb stack mess up the WSL internal DNS.
> You can fix that by editing your `/etc/wsl.conf` and adding the following:
> ```shell
> [network]
> generateResolvConf = false
> ```
> Restart WSL by executing the following on CMD:
> ```shell
> wsl --shutdown
> ```
> Open WSL terminal again and remove the symbolic link from `/etc/resolv.conf`:
> ```shell
> sudo unlink /etc/resolv.conf
> ```
> Create a new `/etc/resolv.conf` file and add the following:
> ```shell
> nameserver 8.8.8.8
> ```
> save the file and you are done.

<a name="kubectl"></a>
## 🐋 k9s (Optional)

Quick install a **k9s** to manage your cluster executing:
```shell
make install-k9s
```

<a name="deploy-orb-kind"></a>
## 🚀  Deploy Orb on Kind

Add `kubernetes.docker.internal` host as `127.0.0.1` address in your hosts file:
```shell
echo "127.0.0.1 kubernetes.docker.internal" | sudo tee -a /etc/hosts
```
> **💡 Note:** This is needed just once

Setup **Orb Charts** dependencies repositories:
```shell
make prepare-helm
```
> **💡 Note:** You just need to run those steps until here once, even if you delete the cluster afterwards.

Use the following command to create the cluster and deploy **Orb**:
```shell
make run
```

Access the **Orb UI** by accessing: https://kubernetes.docker.internal/. The following users are created during the mainflux bootstrap:
E-mail | Password | Role
:--- | :--- | :---
admin@kind.com | pass123456 | Admin

Have fun! 🎉 When you are done, you can delete the cluster by running:
```shell
make kind-delete-cluster
```

## Development flow with Kind


Use the following command to create the empty cluster:
```shell
make kind-create-cluster
```
> **💡 Note:** Now you have and empty kind cluster with minimum necessary to spin up pods


Let's add helm charts for orb:
```shell
make prepare-helm
```
> **💡 Note:** Now your dependencies are configured


Building all orb images:
```shell
make dockers
```

> **💡 Note:** This can take some time

Loading all images into the kind cluster:
```shell
make kind-load-images
```

> **💡 Note:** Your are loading from your local docker registry to kind cluster registry

Load just one image to the kind cluster
```shell
kind load docker-image ns1labs/orb-maestro:0.22.0-088bee14
```

> **💡 Note:** Dont forget to change **kind/values.yaml** manifest to use your image tag


Install orb application:
```shell
make kind-install-orb
```

> **💡 Note:** Now orb was installed properly

If you have any problem to load your new deployment use:
```shell
kubectl rollout restart deployment -n orb
```

<a name="update-service"></a>
## Updating inflight service with recent development
 

If you want to change a service, lets say you added some logs to the fleet service, after commiting the changes, add this
```shell
SERVICE=fleet make build_docker
```
This will build only the docker image of the new service.
After changing you can simply execute
```shell
make kind-upgrade-all
```
