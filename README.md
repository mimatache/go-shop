# go-shop
A simple web shop API written in go


**Things to do**
- Testing - add unit tests
- Comments - add more comments
- Refactor - part of the code is more coupled than I would like

**Usage**
You can run a `make` to lint and build the images

**API**

| Path | Scope |
|------|-------|
| /api/v1/login | Use basic auth to login to the application. This will generate a JWT Token that will be used for subsequent requests |
| /api/v1/logout | Blacklists the current JWT Token so it is no longer usable |
|/api/v1/cart/add | Add a product to the cart. This is a POST requests that expects a message of the form `{"id":1,"quantity":2}` where id is the product ID and quantity is how much you want. This errors out when your cart contains more products of a certain kind that are available in stock. The response contains the current contents of your shopping cart |
|/api/v1/cart/checkout | Attemps to buy the items in your cart. Has a chance to fail. Should block any other people from trying to buy the items|
