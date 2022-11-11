# Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
# SPDX-License-Identifier: Apache-2.0

.script/shutdown.sh
sleep 4
.script/setup.sh
blockMode="-b block"
.script/startup.sh "$blockMode"
cd .mocha || exit
npm run test:awesome