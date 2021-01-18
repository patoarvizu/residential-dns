#!/bin/bash

helm package helm/residential-dns/
version=$(cat helm/residential-dns/Chart.yaml | yaml2json | jq -r '.version')
mv residential-dns-$version.tgz docs/
helm repo index docs --url https://patoarvizu.github.io/residential-dns
helm-docs
mv helm/residential-dns/README.md docs/index.md
git add docs/