# Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
# SPDX-License-Identifier: Apache-2.0

set -x

source ./.script/environment.sh

bash ./.script/shutdown.sh
sleep 5
bash ./.script/setup.sh
bash ./.script/startup.sh

set +x
