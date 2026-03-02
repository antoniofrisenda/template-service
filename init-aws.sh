#!/bin/bash
awslocal s3 mb s3://templates-bucket
awslocal s3 ls

echo "Bucket OK "