package application

const (
	appName     = "mantleNode"
	upgradeName = "v0.3.0"

	authzMsgExec                               = "/cosmos.authz.v1beta1.MsgExec"
	authzMsgGrant                              = "/cosmos.authz.v1beta1.MsgGrant"
	authzMsgRevoke                             = "/cosmos.authz.v1beta1.MsgRevoke"
	bankMsgSend                                = "/cosmos.bank.v1beta1.MsgSend"
	bankMsgMultiSend                           = "/cosmos.bank.v1beta1.MsgMultiSend"
	distributionMsgSetWithdrawAddr             = "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress"
	distributionMsgWithdrawValidatorCommission = "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission"
	distributionMsgFundCommunityPool           = "/cosmos.distribution.v1beta1.MsgFundCommunityPool"
	distributionMsgWithdrawDelegatorReward     = "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"
	feegrantMsgGrantAllowance                  = "/cosmos.feegrant.v1beta1.MsgGrantAllowance"
	feegrantMsgRevokeAllowance                 = "/cosmos.feegrant.v1beta1.MsgRevokeAllowance"
	govMsgVoteWeighted                         = "/cosmos.gov.v1beta1.MsgVoteWeighted"
	govMsgSubmitProposal                       = "/cosmos.gov.v1beta1.MsgSubmitProposal"
	govMsgDeposit                              = "/cosmos.gov.v1beta1.MsgDeposit"
	govMsgVote                                 = "/cosmos.gov.v1beta1.MsgVote"
	stakingMsgEditValidator                    = "/cosmos.staking.v1beta1.MsgEditValidator"
	stakingMsgDelegate                         = "/cosmos.staking.v1beta1.MsgDelegate"
	stakingMsgUndelegate                       = "/cosmos.staking.v1beta1.MsgUndelegate"
	stakingMsgBeginRedelegate                  = "/cosmos.staking.v1beta1.MsgBeginRedelegate"
	stakingMsgCreateValidator                  = "/cosmos.staking.v1beta1.MsgCreateValidator"
	vestingMsgCreateVestingAccount             = "/cosmos.vesting.v1beta1.MsgCreateVestingAccount"
	transferMsgTransfer                        = "/ibc.applications.transfer.v1.MsgTransfer"
)
