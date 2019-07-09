FROM ubuntu

WORKDIR /app

ADD ./main /app

EXPOSE 8080

ENTRYPOINT [ "./main" ]
