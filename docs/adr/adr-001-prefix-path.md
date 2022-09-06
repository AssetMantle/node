# ADR 001: AssetMantle prefix and path 

## Abstract
This is an ADR to discuss the prefix and path to be used for address creation in asset mantle chain.
 
## Context
Addresses generated with assetMantle currently use cosmos as prefix and coin-type 118 to derive bech32 addresses. Should we keep the same with the chain or change prefix to persistence and coin-type 750, or have totally new variables.

## Decision

> This section describes our response to these forces. It is stated in full sentences, with active voice. "We will ..."
> {decision body}

## Consequences

> This section describes the resulting context, after applying the decision. All consequences should be listed here, not just the "positive" ones. A particular decision may have positive, negative, and neutral consequences, but all of them affect the team and project in the future.

### Backwards Compatibility

> All ADRs that introduce backwards incompatibilities must include a section describing these incompatibilities and their severity. The ADR must explain how the author proposes to deal with these incompatibilities. ADR submissions without a sufficient backwards compatibility treatise may be rejected outright.

### Positive

{positive consequences}

### Negative

{negative consequences}

### Neutral

{neutral consequences}

## Further Discussions

While an ADR is in the DRAFT or PROPOSED stage, this section should contain a summary of issues to be solved in future iterations (usually referencing comments from a pull-request discussion).
Later, this section can optionally list ideas or improvements the author or reviewers found during the analysis of this ADR.

## References

- [New hd key path derivation discussion by confio]{https://github.com/confio/cosmos-hd-key-derivation-spec}