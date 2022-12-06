# Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
# SPDX-License-Identifier: Apache-2.0

<<<<<<< HEAD
set -x

source ./.script/environment.sh

bash ./.script/shutdown.sh
sleep 5
bash ./.script/setup.sh
bash ./.script/startup.sh

set +x
=======
.script/shutdown.sh
sleep 4
.script/setup.sh
blockMode="-b block"
.script/startup.sh "$blockMode"
cd .mocha || exit
npm run test:awesome
>>>>>>> parent of b1a33d3 (tests: add shell script)
