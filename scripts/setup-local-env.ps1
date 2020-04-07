#!/env/bin/env pwsh

$ROOTDIR = $PSScriptRoot
$DC_YML = Join-Path $ROOTDIR "docker-compose.minio.yml"

docker-compose -f $DC_YML up -d
$ID=$(docker-compose -f $DC_YML ps -q create_buckets)
$_=$(docker wait $ID)
docker logs $ID