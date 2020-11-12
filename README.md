# go-shop
A simple web shop API written in go



**Usage**
You can run a `make` to lint and build the images

To start the server on port `8080` use:
```sh
PORT=8080 make run-server
```
> :warning: This command blocks. You need to run it in the backgournd of you intend to continue using this terminal session while still having the server running

To start the client on port `8080` use:
```sh
PORT=8080 make run-client
```

If the port is not given, then both client and server start by default on port `9090`.


**API**

| Path | Scope |
|------|-------|
| /api/v1/login | Use basic auth to login to the application. This will generate a JWT Token that will be used for subsequent requests |
| /api/v1/logout | Blacklists the current JWT Token so it is no longer usable |
|/api/v1/cart/add | Add a product to the cart. This is a POST requests that expects a message of the form `{"id":1,"quantity":2}` where id is the product ID and quantity is how much you want. This errors out when your cart contains more products of a certain kind that are available in stock. The response contains the current contents of your shopping cart |
|/api/v1/cart/checkout | Attemps to buy the items in your cart. Has a chance to fail. Should block any other people from trying to buy the items|
