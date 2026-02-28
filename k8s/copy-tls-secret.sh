#!/bin/bash
# wildcard-tls secret'ını devops namespace'inden home namespace'ine kopyala
kubectl get secret wildcard-tls -n devops -o yaml \
  | sed 's/namespace: devops/namespace: home/' \
  | kubectl apply -f -
