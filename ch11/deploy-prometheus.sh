#!/bin/sh

KUBERNETES_SPECS_DIR="kubernetes-specs"
kubectl delete deployment prometheus --context kind-kind --namespace=prometheus
kubectl delete deployment influxdb-exporter --context kind-kind --namespace=prometheus
kubectl apply -f ${KUBERNETES_SPECS_DIR}/prometheus-Namespace.yaml
kubectl apply -f ${KUBERNETES_SPECS_DIR}/prometheus-ServiceAccount.yaml
kubectl apply -f ${KUBERNETES_SPECS_DIR}/prometheus-config-ConfigMap.yaml
kubectl apply -f ${KUBERNETES_SPECS_DIR}/prometheus-Deployment.yaml
kubectl apply -f ${KUBERNETES_SPECS_DIR}/prometheus-Service.yaml
kubectl apply -f ${KUBERNETES_SPECS_DIR}/grafana-Deployment.yaml
kubectl apply -f ${KUBERNETES_SPECS_DIR}/grafana-Service.yaml
kubectl apply -f ${KUBERNETES_SPECS_DIR}/influxdb-generator-Deployment.yaml
echo "----------------------------------------"
echo "Access Prometheus in:\n  - http://localhost:8001/api/v1/namespaces/prometheus/services/prometheus/proxy/\n\n"
echo "To access Grafana run:\n  - kubectl port-forward service/grafana 3000:3000 --context kind-kind --namespace=prometheus\n"
echo "And go to: \n  - http://localhost:3000/"
echo "----------------------------------------"
