.script/shutdown.sh
sleep 4
.script/setup.sh
blockMode="-b block"
.script/startup.sh "$blockMode"
cd .mocha || exit
npm run test:awesome