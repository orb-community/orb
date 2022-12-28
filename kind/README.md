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

Check if you have a **Docker Environment** running by executing:
```shell
docker version
```
If you need help to setup a **Docker Environment**, follow the [steps from here](https://docs.docker.com/engine/install/debian/).

<a name="helm-3"></a>
## ⚓ Helm 3 (Requirement)

Check if you have **Helm 3** installed by executing:
```shell
helm version
```
If you need help to install **Helm 3**, follow the [steps from here](https://helm.sh/docs/intro/install/).

> 🚨 **Warning:** Make sure you have version 3 installed, orb helm charts doesn't officialy support helm 2.

<a name="kubectl"></a>
## 🐋 Kubectl (Requirement)

Check if you have **Kubectl** cmd installed by executing:
```shell
kubectl version --client
```
If you need help to install **Kubectl**, follow the [steps from here](https://kubernetes.io/docs/tasks/tools/).

<a name="install-kind"></a>
## 🚢 Install Kind (Requirement)

Kind is a tool for running local k8s clusters using docker container as nodes.

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

<a name="deploy-orb-kind"></a>
## 🚀  Deploy Orb on Kind

Add `kubernetes.docker.internal` host as `127.0.0.1` address in your hosts file:
```shell
echo "127.0.0.1 kubernetes.docker.internal" | sudo tee -a /etc/hosts
```
> 
Setup **Orb Charts** dependencies repositories:
```shell
make prepare-helm
```
> **💡 Note:** You just need to run those steps until here once, even if you delete the cluster afterwards.

Use the following command to create the cluster and deploy **Orb**:
```shell
make kind-create-all
```

Access the **Orb UI** by accessing: https://kubernetes.docker.internal/. The following users are created during the mainflux bootstrap:
E-mail | Password | Role
:--- | :--- | :---
admin@kind.com | pass123456 | Admin

Have fun! 🎉 When you are done, you can delete the cluster by running:
```shell
make kind-delete-cluster
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
