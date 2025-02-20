// Copyright [2021] - [2025], AssetMantle Pte. Ltd. and the code contributors
// SPDX-License-Identifier: Apache-2.0

package constants

const (
	AuthzMsgExec                        = "/cosmos.authz.v1beta1.MsgExec"
	AuthzMsgGrant                       = "/cosmos.authz.v1beta1.MsgGrant"
	AuthzMsgRevoke                      = "/cosmos.authz.v1beta1.MsgRevoke"
	BankMsgSend                         = "/cosmos.bank.v1beta1.MsgSend"
	BankMsgMultiSend                    = "/cosmos.bank.v1beta1.MsgMultiSend"
	DistrMsgSetWithdrawAddr             = "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress"
	DistrMsgWithdrawValidatorCommission = "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission"
	DistrMsgFundCommunityPool           = "/cosmos.distribution.v1beta1.MsgFundCommunityPool"
	DistrMsgWithdrawDelegatorReward     = "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"
	FeegrantMsgGrantAllowance           = "/cosmos.feegrant.v1beta1.MsgGrantAllowance"
	FeegrantMsgRevokeAllowance          = "/cosmos.feegrant.v1beta1.MsgRevokeAllowance"
	GovMsgVoteWeighted                  = "/cosmos.gov.v1beta1.MsgVoteWeighted"
	GovMsgSubmitProposal                = "/cosmos.gov.v1beta1.MsgSubmitProposal"
	GovMsgDeposit                       = "/cosmos.gov.v1beta1.MsgDeposit"
	GovMsgVote                          = "/cosmos.gov.v1beta1.MsgVote"
	StakingMsgEditValidator             = "/cosmos.staking.v1beta1.MsgEditValidator"
	StakingMsgDelegate                  = "/cosmos.staking.v1beta1.MsgDelegate"
	StakingMsgUndelegate                = "/cosmos.staking.v1beta1.MsgUndelegate"
	StakingMsgBeginRedelegate           = "/cosmos.staking.v1beta1.MsgBeginRedelegate"
	StakingMsgCreateValidator           = "/cosmos.staking.v1beta1.MsgCreateValidator"
	VestingMsgCreateVestingAccount      = "/cosmos.vesting.v1beta1.MsgCreateVestingAccount"
	TransferMsgTransfer                 = "/ibc.applications.transfer.v1.MsgTransfer"
	LiquidityMsgCreatePool              = "/tendermint.liquidity.v1beta1.MsgCreatePool"
	LiquidityMsgSwapWithinBatch         = "/tendermint.liquidity.v1beta1.MsgSwapWithinBatch"
	LiquidityMsgDepositWithinBatch      = "/tendermint.liquidity.v1beta1.MsgDepositWithinBatch"
	LiquidityMsgWithdrawWithinBatch     = "/tendermint.liquidity.v1beta1.MsgWithdrawWithinBatch"
)
