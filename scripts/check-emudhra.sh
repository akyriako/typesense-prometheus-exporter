docker run --rm \
  --entrypoint sh \
  typesense-prometheus-exporter:ca-remediation \
  -c '
    ls /usr/share/ca-certificates/mozilla/ | grep -Ei "emSign_ECC_Root_CA_-_G3|Comodo_AAA" || true
  '