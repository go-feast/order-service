FROM ubuntu:latest
LABEL authors="viacheslav"

ENTRYPOINT ["top", "-b"]