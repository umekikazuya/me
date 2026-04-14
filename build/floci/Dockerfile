FROM hectorvent/floci:latest-aws

ENV FLOCI_HOSTNAME=floci
ENV FLOCI_STORAGE_MODE=persistent

COPY init/01.sh /etc/floci/init/start.d/01.sh

# ヘルスチェック
HEALTHCHECK --interval=5s --timeout=5s --retries=12 \
  CMD aws --endpoint-url http://localhost:4566 dynamodb describe-table --table-name "me." >/dev/null 2>&1