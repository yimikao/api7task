# api7task



```yaml
apiVersion: apisix.apache.org/v2
kind: ApisixRoute
metadata:
  name: api-routes
spec:
  http:
    - name: route-1
      match:
        hosts:
          - api7.task
        paths:
          - /home
      backends:
        - serviceName: backend-api
          servicePort: 8080
      plugins:
        - name: limit-count
          enable: true
          config:
            count: 10
            time_window: 10
        - name: key-auth
          enable: true
          config:
            key: "auth-one"

```
To run
```sh
eksctl create cluster \
    --name api7task-cluster \
	--version 1.25 \
	--node-type t3.large \
	--nodes 4 \


eksctl utils associate-iam-oidc-provider --cluster=api7task-cluster --region us-west-2 --approve


eksctl create iamserviceaccount \
	--name ebs-csi-controller-sa \
	--namespace kube-system \
	--cluster api7task-cluster \
	--attach-policy-arn arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy \
	--approve \
	--role-only \
	--role-name AmazonEKS_EBS_CSI_DriverRole


eksctl create addon --name aws-ebs-csi-driver --cluster api7task-cluster --service-account-role-arn arn:aws:iam::<account-id>:role/AmazonEKS_EBS_CSI_DriverRole --force


kubectl run backend-api --image yimikao/api7task:v1 --port 8080 -- 8080

kubectl expose pod backend-api --port 8080

set ADMIN_API_VERSION v3

helm install apisix apisix/apisix \
	--set gateway.type=LoadBalancer \
	--set ingress-controller.enabled=true \
	--create-namespace \
	--namespace ingress-apisix \
	--set ingress-controller.config.apisix.serviceNamespace=ingress-apisix \
	--set ingress-controller.config.apisix.adminAPIVersion=$ADMIN_API_VERSION

kubectl apply -f manifests/route.yaml

kubectl get --namespace ingress-apisix svc -w apisix-gateway

export SERVICE_IP=$(kubectl get svc --namespace ingress-apisix apisix-gateway --template "{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}")

echo http://$SERVICE_IP:80

kubectl exec -n ingress-apisix deploy/apisix -- curl http://127.0.0.1:9180/apisix/admin/consumers -X PUT -d '
	{
		"username": "yinka",
		"plugins": {
			"key-auth": {
				"key": "auth-one"
			}
		}
	}'


curl http://a0aac4250495046a6a989418b0e3cf58-1370549474.us-west-2.elb.amazonaws.com:80/home -H 'host:api7.task' -H 'apikey:auth-one'


```

