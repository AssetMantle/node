# Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
# SPDX-License-Identifier: Apache-2.0

set -x

source environment.sh

bash shutdown.sh
sleep 5
bash setup.sh
bash startup.sh

set +x