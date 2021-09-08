# Orb Helm Chart

Helm Chart for the Orb Observability Platform

## Prerequisites

- Helm v3

## Configuration

## Instructions

This guide assumes installation into name space `orb`. It requires a HOSTNAME you have DNS control over. It uses Let's
Encrypt for TLS certification management.

* cd to working directory `charts/orb`
* Add helm repos for dependencies

```
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm dependency update
```

* Create `orb` namespace

```
kubectl create namespace orb
```

* Create JWT signing key secret

```
kubectl create secret generic orb-auth-service --from-literal=jwtSecret=MY_SECRET -n orb
```

* Create admin user secrets

```
kubectl create secret generic orb-user-service --from-literal=adminEmail=user@example.com --from-literal=adminPassword=12345678 -n orb
```

* Deploy [ingres-nginx helm](https://kubernetes.github.io/ingress-nginx/deploy/#using-helm) (to default namespace) with
  tcp config map configured from helm for 8883 (MQTTS). Note you need to reference both namespace and helm release name
  here!

```
helm install --set tcp.8883=orb/my-orb-nginx-internal:8883 ingress-nginx ingress-nginx/ingress-nginx
```

* Wait for an external IP to be available

```
kubectl --namespace default get services -o wide -w ingress-nginx-controller
```

* Choose a HOSTNAME, then point an A record for this hostname to the external IP
* Deploy [cert manager helm](https://cert-manager.io/docs/installation/helm/)
  to [secure nginx ingress](https://cert-manager.io/docs/tutorials/acme/ingress/)

```
helm install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --version v1.5.3 --set installCRDs=true
```

* Create Issuer CRDs (in the `orb` namespace!)
  * `cp issuers/production-issuer-tpt.yaml issuers/production-issuer.yaml`
  * edit `issuers/production-issuer.yaml` and change `spec.acme.email` to a real email address
  * `kubectl create -f issuers/production-issuer.yaml -n orb`

* Install orb. Replace `my-orb` with your helm release name.

```
helm install --set ingress.hostname=HOSTNAME -n orb my-orb .
```

