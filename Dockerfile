FROM alpine

WORKDIR /app

COPY getsrc-linux-amd64 /app/
COPY ./tmpl /app/tmpl/
COPY ./static /app/static/

ENTRYPOINT ["/app/getsrc-linux-amd64"]