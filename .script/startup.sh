assetNode start >~/.AssetMantle/Node/log &
sleep 10
assetClient rest-server --chain-id test $1 $2 >~/.AssetMantle/Client/log &
echo "
Node and Client started up. For logs:
tail -f ~/.AssetMantle/Node/log
tail -f ~/.AssetMantle/Client/log
"
