#!/bin/sh

KUBERNETES_SPECS_DIR="kubernetes-specs"

kubectl apply -f ${KUBERNETES_SPECS_DIR}/prometheus-Namespace.yaml
kubectl apply -f ${KUBERNETES_SPECS_DIR}/prometheus-ServiceAccount.yaml
kubectl apply -f ${KUBERNETES_SPECS_DIR}/prometheus-ClusterRole.yaml
kubectl apply -f ${KUBERNETES_SPECS_DIR}/prometheus-ClusterRoleBinding.yaml
kubectl apply -f ${KUBERNETES_SPECS_DIR}/prometheus-config-ConfigMap.yaml
kubectl apply -f ${KUBERNETES_SPECS_DIR}/prometheus-Deployment.yaml
kubectl apply -f ${KUBERNETES_SPECS_DIR}/prometheus-Service.yaml

