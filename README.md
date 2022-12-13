# timestamps
A JSON/HTTP service, in golang, that returns the matching timestamps of a periodic task.

#### Run without Docker:
In the root directory of the project, type the following to fire up the server:
```
chmod +x run.sh && ./run.sh 
```
 (don't forget to add port to listen to, for example : `./run.sh 8080`)

#### Run with Docker:
In the root directory of the project (where the dockerfile exists):
```
docker build --tag docker-inaccess .
```

To view local images (docker-inaccess image should be here):
```
docker images
```
To run the docker-inaccess image:
```
docker run -p 8080:8080 docker-inaccess 8080
```
It's handy to remember that in the -p 8080:8080 command line option, the order for this operation is `-p HOST_PORT:CONTAINER_PORT`

In case you want to delete the image:
```
docker image rm -f docker-inaccess
```
