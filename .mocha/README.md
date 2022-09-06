## Pre-requisites

Install npm and nodeJS

```
cd ~
curl -sL https://deb.nodesource.com/setup_10.x -o nodesource_setup.sh
sudo bash nodesource_setup.sh
sudo apt install nodejs
```

To check which version of Node.js you have installed after these initial steps, type:

```
node -v
```

For more information, visit https://www.digitalocean.com/community/tutorials/how-to-install-node-js-on-ubuntu-18-04.


* * *

## Installation

```
cd assetMantle/.mocha
npm install
```

## Testing

To test, go to the assetMantle/.mocha folder and type :

```
./.script/resetChainAndRunMocha.sh
```

or

```
npm run test:awesome
```

This will reset the chain and run mocha tests

NOTE: If any error comes which says: Error: Cannot find module 'xxx' then run "npm install xxx --save"

***

## Documentation

For more information, visit https://autom8able.com.

***

## Report

[mochawesome] Report JSON saved to assetMantle/.mocha/mochawesome-report/mochawesome.json

[mochawesome] Report HTML saved to assetMantle/.mocha/mochawesome-report/mochawesome.html


* * *