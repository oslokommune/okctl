# 0003. Emoji git log

<!-- 
    Date representing when the RFC was submittet for review 

    Example:
    Date: 08.02.2021
-->
Date: 15.02.2021

## Context

<!--
    Short description of why this decision is needed

    Example:
    We want to better document our decisions
-->

Contributors should be able to easily create meaningful commit messages.

I propose to use the following as a style for our commit messages:

Inspired by [Emoji-log](https://ahmadawais.com/emoji-log/).

The format looks like this:
`[emoji] [change as present tense imperative action] [issue ref]`

Imperative meaning; write it like you're giving an order.

1. ✅ Add/create things. Be it functionality, tests, documentation. Examples:
   - `✅ Add nightly test #1337`
   - `✅ Add declarative cluster KM98`
   - `✅ Add declarative cluster user documentation #400`
   - `✅ Add observability stack KM88`
2. 👌 Improve things. Cleaning out some code, refactoring,
   - `👌 Improve test coverage for create cluster #132`
   - `👌 Improve readability in apply cluster code #144`
   - `👌 Improve error handling for delete cluster KM37`
   - `👌 Bump git sdk to 1.24`
3. 🐛 Fix things. Bugs, typos
   - `🐛 Fix delete cluster not finding ArgoCD domain KM89`
   - `🐛 Fix rogue pointer print line in apply cluster KM44`
   - `🐛 Fix missing dependency logrus #68`
   - `🐛 Fix confusing error message when failing to create VPC #12`

## Decision

<!--
    Short summary of the decision

    Example:
    We've decided to go for AlphaGov's system as explained [here](https://github.com/alphagov/govuk-aws/blob/24d1ea513e58ee938043d71d09815a51229067bf/docs/architecture/decisions/0001-record-architecture-decisions.md)
-->

Until the need for uniform commit messages arise, the contributors will use what they see fit as long as they adhere to
[this](https://github.com/oslokommune/okctl/blob/master/docs/rfc/0002-git-commit-messages.md) RFC.

This RFC and [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/) are the two most common styles today.

## Consequences

<!--
    Short description of positive and negative consequences of the decision

    Example:
    The decisions we deem relevant will be documented in this format here in this folder.
-->

### Positive:
- Uniform way to read and write the git log
- Easily determine the purpose of the change (It's much easier skimming for an emoji rather than text)
- Fancy looking git log
- A simple system to learn and use

### Negative
- Probably not easy to automatically generate semantic versions from
