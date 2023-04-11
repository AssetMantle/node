// Copyright [2021] - [2022], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package queuing

import (
	"github.com/Shopify/sarama"
)

// kafkaAdmin : is admin to create topics
func kafkaAdmin(kafkaNodes []string) sarama.ClusterAdmin {
	config := sarama.NewConfig()
	config.Version = sarama.V0_11_0_0 // hardcoded

	admin, err := sarama.NewClusterAdmin(kafkaNodes, config)
	if err != nil {
		panic(err)
	}

	return admin
}
