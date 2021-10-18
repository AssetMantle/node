assetNode start >~/.assetNode/log &
sleep 10
assetClient rest-server -b sync --unsafe-cors --generate-only  --chain-id test >~/.assetClient/log &
echo "
Node and Client started up. For logs:
tail -f ~/.assetNode/log
tail -f ~/.assetClient/log
"
