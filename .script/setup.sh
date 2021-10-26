make all

rm -rf ~/.assetNode
rm -rf ~/.assetClient

mkdir ~/.assetNode
mkdir ~/.assetClient

assetNode init --chain-id test test
assetClient keys add test <<<"n" >/dev/null 2>&1 0>&1
assetClient keys add test --recover<<<"y
wage thunder live sense resemble foil apple course spin horse glass mansion midnight laundry acoustic rhythm loan scale talent push green direct brick please"
assetNode add-genesis-account test 100000000000000umantle
assetNode gentx --name test --amount 1000000000umantle
assetNode collect-gentxs
