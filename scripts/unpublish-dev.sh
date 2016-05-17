#!/usr/bin/env bash

release_url='https://api.github.com/repos/dirkraft/ash/releases/3238640'
delete_url='https://api.github.com/repos/dirkraft/ash/releases/assets'
asset_ids=$(curl --silent "${release_url}" | jq '.assets[].id')

for asset_id in ${asset_ids} ; do

  curl -XDELETE --silent \
    --header "Authorization: token ${GITHUB_TOKEN}" \
    "${delete_url}/${asset_id}"

done
