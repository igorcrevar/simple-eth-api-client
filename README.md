Modified code from https://hackernoon.com/create-an-api-to-interact-with-ethereum-blockchain-using-golang-part-1-sqf3z7z
different route handling, some changes inside modules, added chain api call

Use node from https://dashboard.alchemyapi.io/ (Bad boy has one)
 - url would be something like https://eth-rinkeby.alchemyapi.io/v2/YOUR_API_KEY
or install https://www.trufflesuite.com/ganache (by default if no arguments provided, program will use ganache node)

- ...you can also specify private key for each transaction (or you can specify `privateKey` inside `Body` each time)
