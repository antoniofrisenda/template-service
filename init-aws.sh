#!/bin/bash
set -e

awslocal s3 mb s3://document-bucket
awslocal s3 ls

echo "Bucket OK "