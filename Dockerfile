FROM golang:1.5.2-onbuild

EXPOSE 8080
WORKDIR /go/src/app
CMD ["./run_docker.sh"]