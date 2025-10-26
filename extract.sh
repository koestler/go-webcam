id=$(docker create docker.io/library/go-webcam:latest)
docker cp $id:/go-webcam .
docker rm -v $id

