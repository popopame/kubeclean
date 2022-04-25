# Kubeclean

DISCLAIMER: This repository is for learning purposes, the original project is kubectl-neat


Kubeclean is golang CLI tool that is used to clean kubernetes manifest.

When a Kubernetes resource is created, the cluster will add fields in the resource definition (status, metadata, ServiceAccount,...).

Theses fields may block you from easily reading the object configuration and if you try to apply the object as is, Kubernetes will return an error

