[![Acceptance Tests](https://github.com/SatyaLens/source-score/actions/workflows/acceptance-tests.yml/badge.svg)](https://github.com/SatyaLens/source-score/actions/workflows/acceptance-tests.yml)   

# Source Score

## Overview

Source Score is a microservice that rates information sources based on the validity of their claims. The system evaluates sources by analyzing claims made by those sources and the supporting or refuting evidence (proofs) for each claim. Sources receive a score between 0 and 1, calculated as the ratio of valid claims to total verified claims.

## Data Models

### Source

A source represents an information provider (e.g., news outlet, research institution, blog) that makes claims. Each source has:

- **name**: Display name of the source
- **uri**: Unique HTTPS URL identifying the source
- **uriDigest**: SHA-256 hash of the URI, used as the primary key
- **summary**: Brief description of the source
- **tags**: Comma-separated categorization tags
- **score**: Calculated credibility score (0.0 to 1.0) based on the ratio of valid to total verified claims

### Claim

A claim represents a statement or assertion made by a source. Each claim has:

- **sourceUriDigest**: Reference to the parent source
- **title**: Short title of the claim
- **summary**: Detailed description of the claim
- **uri**: Unique HTTPS URL identifying the claim
- **uriDigest**: SHA-256 hash of the URI, used as the primary key
- **checked**: Boolean indicating whether the claim has been verified
- **validity**: Boolean indicating whether the claim is valid (true) or invalid (false) based on proof analysis

### Proof

A proof represents evidence that either supports or refutes a claim. Each proof has:

- **claimUriDigest**: Reference to the claim being evaluated
- **uri**: Unique HTTPS URL identifying the proof source
- **uriDigest**: SHA-256 hash of the URI, used as the primary key
- **supportsClaim**: Boolean indicating whether this proof supports (true) or refutes (false) the claim
- **reviewedBy**: Identifier of the reviewer who evaluated this proof

## How It Works

1. **Create Sources**: Register information sources in the system
2. **Add Claims**: Associate claims with their respective sources
3. **Submit Proofs**: Add evidence that supports or refutes each claim
4. **Verify Claims**: Trigger the verification process that analyzes all proofs for each claim and determines validity (more supporting proofs = valid, more refuting proofs = invalid)
5. **Update Scores**: Calculate source scores based on the ratio of valid claims to total verified claims
