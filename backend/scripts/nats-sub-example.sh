# Note: replace the user (and subject) with the workspace id
nats --server nats://127.0.0.1:4222 \
     --user "b60e60b0-4891-44a5-ab13-e1a25c70c1fa" \
     --password "$(bash scripts/get-token.sh)" \
     sub "workspaces.b60e60b0-4891-44a5-ab13-e1a25c70c1fa.>"
