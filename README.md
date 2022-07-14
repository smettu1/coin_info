# Setup
Setup consists of two pods, the first one is go app, and the second is redis.

# Design
Simple ETL application which receives coin list as POST requests and gets the meta information 
related to coin and saves in SQL db.
Go process receives the array of coin requests and gets the response from 
coingeko API server and generates the output from the API response.
To make ETL truly scalable go process needs to be made stateless and highly concurrent.

In case coingeko throttles API's they need to be retired and in case of failure
these request objects need to be stored in dead letterQ.
Once there's a successful response then the go process generates taskId.

# Improvements
Create session objects and pass them around in coroutine in a thread-safe way.

# Usage
docker-compose build builds docker images
docker-compose up runs the server.

ttu@shravans-mbp ~ % curl -X POST http://localhost:3000/coins -H 'Content-Type: application/json' -d '["airswap"]'
ttu@shravans-mbp ~ % curl http://localhost:3000/output
[{"Id":"airswap","exchanges":["gdax","binance","xt","gate","latoken","bitrue","uniswap","coinex","huobi","bkex","okex","hotbit","huobi","gate","nice_hash","balancer_v1","uniswap_v2","hoo"],"TaskRun":0}]
ttu@shravans-mbp ~ % curl -X POST http://localhost:3000/coins -H 'Content-Type: application/json' -d '["aitheon"]'
ttu@shravans-mbp ~ % curl -X POST http://localhost:3000/coins -H 'Content-Type: application/json' -d '["aiwork"]'
ttu@shravans-mbp ~ % curl http://localhost:3000/output
[{"Id":"airswap","exchanges":["gdax","binance","xt","gate","latoken","bitrue","uniswap","coinex","huobi","bkex","okex","hotbit","huobi","gate","nice_hash","balancer_v1","uniswap_v2","hoo"],"TaskRun":0},{"Id":"aiwork","exchanges":["bithumb","mxc","bkex"],"TaskRun":1}]


