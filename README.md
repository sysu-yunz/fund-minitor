# Fund Monitor

## Deploy

### vercel

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fsysu-yunz%2Ffund-minitor&env=MGO_PWD,BOT_TOKEN&envDescription=You%20need%20a%20Telegram%20bot%20and%20MongoDB%20Atlas%20account%20to%20get%20ready.&project-name=my-tg-bot)


### heroku


A telegram bot that provide global stock market index and crypto price.

Using webhook deployed on vercel.com and local dev with telegram polling mode.

## Subscribe

Users can subscribe quotes they like. Send command /quotes, bot will reply all of your subscribed quotes.

## Check

Send quote name to check the latest price. Bot will find your favorite by your history.

## Buy

Buy command with price and amount, bot will build a portfolio for you, you can change it by buy/sell command.

## Data

### Fund

Fund data source from http://api.fund.eastmoney.com

### Crypto


### Commodities

### Crypto

### Indices

### Stocks

Stock data source from xueqiu.com

Stock list in https://xueqiu.com/hq#exchange=US&firstName=3&secondName=3_0

Stock detail in https://xueqiu.com/S/SZ000002

## Statement

Data only used for research purpose.

