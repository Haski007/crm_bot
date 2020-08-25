# TELEGRAM BOT (Pre-Alpha v.0.32) -pocket crm for small business

## There are two types of users:
* users - hired workers. They can add to database new purchases and create new products.
* admins - founderds or cofounders of business. They can remove users, add new products. They receive statistics every day at 23:00 by Europe/Kiev.
## Statistics:
* All sold products for whole day
* Quantity of sold products
* Total amount
## History:
* All purchases that were made for today
* The ability to delete any purchase
## Storage:
* The ability to track quantity of all of your products at stock at this moment
* Bot Warns you if quantity of product is less than 10
## Future features:
* Statistics for whole month
* Statistics which product sells best and which sells least
### To register you need to write your boss or author of bot https://t.me/pdemian

![](/assets/images/adding_purchase.gif)
![](/assets/images/getting_statistics.gif)

## Bot hosted on Google Cloud Platform. Database - mongoDB.

`docker-compose up -d --build` - to build docker-compose with mongoDB
