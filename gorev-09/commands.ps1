docker compose build
docker run --rm -v "${PWD}:/src" anchore/syft:latest dir:/src -o cyclonedx-json > gorev-09/sbom.json
docker run --rm -v "${PWD}:/src" aquasec/trivy:latest fs --severity CRITICAL,HIGH /src
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy:latest image --severity CRITICAL,HIGH secscan-backend:latest
docker run --rm -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy:latest image --severity CRITICAL,HIGH secscan-frontend:latest
