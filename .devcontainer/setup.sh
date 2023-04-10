#!/usr/bin/env bash
mkdir -p $HOME/.kube
sudo cp -r /usr/local/share/kube-localhost/* $HOME/.kube
sudo chown -R $(id -u) $HOME/.kube

KUBECTL_VERSION=$(curl -L -s https://dl.k8s.io/release/stable.txt)
curl -LO "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl"
chmod +x kubectl
sudo mv kubectl "/usr/local/bin/"

curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh
rm get_helm.sh