FROM alpine:3.7

COPY data /etc/data

COPY shop /usr/bin/shop

ENV PORT=9090

CMD shop -port ${PORT} -users /etc/data/users.json -products /etc/data/products.json 
