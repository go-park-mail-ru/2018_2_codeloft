#!/bin/bash
git checkout testing
git pull
docker build --build-arg USERNAME=codeloft --build-arg PASSWORD=1codeloft1 -t codeloftbuild .
docker stop codeloft
docker rm codeloft
docker rmi codeloft
docker tag codeloftbuild codeloft
docker rmi codeloftbuild
docker run -d -p 3000:8080 --name codeloft -t codeloft
