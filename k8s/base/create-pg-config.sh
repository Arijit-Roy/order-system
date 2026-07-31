#!/bin/bash
kubectl create configmap postgres-init \
  --from-file=init.sql=k8s/postgres/sql/init.sql \
  --dry-run=client \
  -o yaml > k8s/postgres/postgres-init.yaml