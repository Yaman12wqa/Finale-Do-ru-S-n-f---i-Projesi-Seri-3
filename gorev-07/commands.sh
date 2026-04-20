#!/usr/bin/env bash
set -euo pipefail

docker compose up -d
nmap -sV -p 3000,8080 localhost
docker run --rm -t -v "$(pwd)/gorev-07:/zap/wrk" ghcr.io/zaproxy/zaproxy:stable zap-baseline.py -t http://host.docker.internal:3000 -r zap-report.html
