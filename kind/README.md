# Orb.live Helm Chart
Helm Chart for the SaaS Orb Observability Platform

## Prerequisites
- Helm v3
- **openssl** command line tool
- **jq** command line tool
- **aws** command line tool
- **dockerhub** private repos (ns1labs/orb-ui-live-{{environment}}) permissions 
- **aws resources** (eks, s3, dynamodb and sts) permissions

## Instructions
This guide assumes installation into name space `orb`. It requires a HOSTNAME you have DNS control over.
It uses Let's Encrypt for TLS certification management.

### Preparing your environment
1. Enter the orb.live chart directory:
   ```shell
   cd charts/orb-live
   ```
   
2. Add helm repos for dependencies
   ```shell
   helm repo add jaegertracing https://jaegertracing.github.io/helm-charts
   helm repo add bitnami https://charts.bitnami.com/bitnami
   helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
   helm repo add jetstack https://charts.jetstack.io
   helm repo update
   helm dependency update
   ```

3. Get credentials to the K8s cluster and SMTP server from terraform outputs
    ```shell
    ENVIRONMENT=# dev, stage or prod
    export KUBECONFIG=$(mktemp -t eks-${ENVIRONMENT}.conf.XXXX) &&
    cd ../../terraform/environments/${ENVIRONMENT}/orb-app &&
    SMTP_CREDS=$(terraform output --json smtp_credentials) &&
    SMTP_HOST=$(echo $SMTP_CREDS | jq .server) &&
    SMTP_USER=$(echo $SMTP_CREDS | jq .username) &&
    SMTP_PASS=$(echo $SMTP_CREDS | jq .password) &&
    cd ../k8s-cluster &&
    aws eks --region $(terraform output --raw region) update-kubeconfig \
      --name $(terraform output --raw cluster_name) \
      --role-arn $(terraform output --raw management_role_arn) --dry-run > $KUBECONFIG &&
    cd ../../../../charts/orb-live/
    ```

### Deploying Orb app to the Kubernetes cluster
1. Create the following secrets in the `orb` namespace:
    * Create the `orb` namespace
     ```shell
     kubectl create namespace orb
     ```
    * Create the JWT signing key secret
     ```shell
     JWT_SECRET=$(openssl rand -base64 27)
     kubectl create secret -n orb generic orb-auth-service --from-literal=jwtSecret=${JWT_SECRET} 
     ``` 
     
   * Create keto dns secret
   ```
   kubectl create secret generic orb-keto-dsn --from-literal=dsn='postgres://postgres:password@db.host.com:5432/keto' -n orb
   ```

     
    * Create sinks encryption password
     ```
     kubectl create secret generic orb-sinks-encryption-key --from-literal=key=mainflux-orb -n orb-live
     ```
    * Create the admin user credentials secret
     ```shell
     ADMIN_EMAIL=# something like admin@dev.orb.live
     ADMIN_PASSWORD=$(openssl rand -base64 18)
     kubectl create secret -n orb generic orb-user-service \
       --from-literal=adminEmail=${ADMIN_EMAIL} --from-literal=adminPassword=${ADMIN_PASSWORD}
     ```    
    * Create the Dockerhub credentials for custom UI image.
    Make sure you're logged in to dockerhub, [instructions here](https://docs.docker.com/engine/reference/commandline/login/).
     ```shell
     kubectl create secret -n orb docker-registry orb-ui-regcred \
      --from-file=.dockerconfigjson=$HOME/.docker/config.json
     ```
    
2. Create the ingress Load balancer:
    * Install nginx ingress
      ```shell
      helm install -n default --set tcp.8883=orb/my-orb-nginx-internal:8883 \
        ingress-nginx ingress-nginx/ingress-nginx
      ```
    * Add a DNS alias record having the HOSTNAME point to the nginx-ingress public address.
      The following command helps to get the public address:
      ```shell
      kubectl -n default get services -o json ingress-nginx-controller \
        | jq -r .status.loadBalancer.ingress[0].hostname
      ```
    * Install cert-manager to enable HTTPS
      ```shell
      helm install -n cert-manager --create-namespace --set installCRDs=true \
        cert-manager jetstack/cert-manager --version v1.5.3
      
      kubectl create -n orb -f ../../certs/production-issuer.yaml
      ```
   
3. Deploy Orb to kubernetes, passing in the correct values file and the SMTP credentials from terraform deployment:
    ```shell
    helm install -n orb --values="values-${ENVIRONMENT}.yaml" --set orb.smtp.host=${SMTP_HOST} \
      --set orb.smtp.username=${SMTP_USER} --set orb.smtp.password=${SMTP_PASS} my-orb .
    ```

### Deleting Orb app from the Kubernetes cluster
1. Remove all created namespaces
    ```shell
    kubectl delete namespace orb &&
    kubectl delete namespace cert-manager
    ```
2. Remove the ingress load-balancer
    ```shell
    helm uninstall ingress-nginx
    ```


### Installing kubernetes-dashboard
This step need to be executed only when needed. The orb.live deploy doesn't depends on it. You'll need to execute steps 1 and 3 of "Preparing your environment" section above before doing the steps below.
1. Add dependencies
   ```shell
   helm repo add kubernetes-dashboard https://kubernetes.github.io/dashboard/
   ```
2. Create kubernetes-dashboard namespace
   ```shell
   kubectl create namespace kubernetes-dashboard
   ```
3. Install kubernetes-dashboard
   ```shell
   helm install -n kubernetes-dashboard kubernetes-dashboard kubernetes-dashboard/kubernetes-dashboard
   ```
4. Apply rbac model
   ```shell
   kubectl apply -R -f manifests/rbac/users.d
   kubectl apply -R -f manifests/rbac/roles.d
   kubectl apply -R -f manifests/rbac/users-to-role-binding.d
   ```
5. To grant access to kubernetes dashboard, you can use one of these two strategies:
   1. Port-forward. This whill redirect the kubernetes-dashboard pod port inside of the kubernetes cluster to your local machine.
      ```shell
      export POD_NAME=$(kubectl get pods -n kubernetes-dashboard -l "app.kubernetes.io/name=kubernetes-dashboard,app.kubernetes.io/instance=kubernetes-dashboard" -o jsonpath="{.items[0].metadata.name}")
      echo https://127.0.0.1:8443/
      kubectl -n kubernetes-dashboard port-forward $POD_NAME 8443:8443  > /dev/null 2>&1 &
      ```
   2. Ingress policy. This will create a policy inside of the ingress service to receive external connection. Note that to use this strategy, the NS1 DNS need to be configured with kubernetes.orb.live appointing to the ingress hostname.
      ```shell
      kubectl apply -f manifests/ingress/kubernetes-dashboard.yaml
      ```
6. Get your k8s-dash access token
   ```shell
   NAME=eric \
   kubectl -n kubernetes-dashboard get secret $(kubectl -n kubernetes-dashboard get secret sa/${NAME}-user -o jsonpath="{.secrets[0].name}") -o go-template="{{.data.token | base64decode}}"
   ```
