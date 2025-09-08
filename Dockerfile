FROM ubuntu:latest
LABEL authors="angro"

ENTRYPOINT ["top", "-b"]