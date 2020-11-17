FROM alpine:3.7

COPY data /etc/data

COPY shop /usr/bin/shop

ENV PORT=9090

ENV USERS="/etc/data/users.json"

ENV PRODUCTS="/etc/data/products.json"

CMD shop -port ${PORT} -users ${USERS} -products ${PRODUCTS}
