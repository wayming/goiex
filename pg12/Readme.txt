Build image:
docker build --progress=plain -t goiex-postgres-centos .

Run container:
docker kill pg; docker system prune -f; docker run --privileged -d --name pg goiex-postgres-centos

Attach the container:
docker exec -it pg bash
