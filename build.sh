#!/usr/bin/env bash
VERSION="$1"
if [[ "${VERSION}" == "" ]]; then
  echo "Version parameter needed"
  exit
fi

docker build . -t mortbury.azurecr.io/frontdoor-ingress:${VERSION} --build-arg="version=${VERSION}"
docker push mortbury.azurecr.io/frontdoor-ingress:${VERSION}
helm package ./helm/frontdoor-ingress --app-version=${VERSION} --version=${VERSION}
helm upgrade --wait --install --namespace ingress --create-namespace --set="image.tag=${VERSION}" fdingress ./frontdoor-ingress-${VERSION}.tgz
rm frontdoor-ingress-${VERSION}.tgz