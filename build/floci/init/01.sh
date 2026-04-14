#!/bin/sh
set -eu

aws dynamodb create-table \
  --endpoint-url http://localhost:4566 \
  --table-name me. \
  --billing-mode PAY_PER_REQUEST \
  --attribute-definitions \
  AttributeName=PK,AttributeType=S \
  AttributeName=SK,AttributeType=S \
  AttributeName=GSI1PK,AttributeType=S \
  AttributeName=GSI1SK,AttributeType=S \
  AttributeName=GSI2PK,AttributeType=S \
  AttributeName=GSI2SK,AttributeType=S \
  AttributeName=GSI3PK,AttributeType=S \
  AttributeName=GSI3SK,AttributeType=S \
  AttributeName=GSI_EMAIL_PK,AttributeType=S \
  --key-schema \
  AttributeName=PK,KeyType=HASH \
  AttributeName=SK,KeyType=RANGE \
  --global-secondary-indexes \
  'IndexName=GSI1,KeySchema=[{AttributeName=GSI1PK,KeyType=HASH},{AttributeName=GSI1SK,KeyType=RANGE}],Projection={ProjectionType=ALL}' \
  'IndexName=GSI2,KeySchema=[{AttributeName=GSI2PK,KeyType=HASH},{AttributeName=GSI2SK,KeyType=RANGE}],Projection={ProjectionType=ALL}' \
  'IndexName=GSI3,KeySchema=[{AttributeName=GSI3PK,KeyType=HASH},{AttributeName=GSI3SK,KeyType=RANGE}],Projection={ProjectionType=ALL}' \
  'IndexName=GSI_EMAIL,KeySchema=[{AttributeName=GSI_EMAIL_PK,KeyType=HASH},{AttributeName=SK,KeyType=RANGE}],Projection={ProjectionType=ALL}'
