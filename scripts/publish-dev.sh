#!/usr/bin/env bash

upload_url='https://uploads.github.com/repos/dirkraft/ash/releases/3238640/assets'
asset_file=$1
asset_name=$(basename $asset_file)

curl --silent \
  --header "Authorization: token ${GITHUB_TOKEN}" \
  --header "Content-Type: application/octet-stream" \
  "${upload_url}?name=${asset_name}" \
  --data-binary "@${asset_file}" | jq .
