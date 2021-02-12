# 0. Request for comments

<!-- 
    Date representing when the RFC was submittet for review 

    Example:
    Date: 08.02.2021
-->
Date: 12.02.2021

## Status
<!--
    Uncomment relevant status
-->

Pending

<!--
Accepted
Rejected
-->

## Context

<!--
    Short description of why this decision is needed

    Example:
    We want to better document our decisions
-->

We want better traceability for our git commits. We want this by prefixing commit messages by a github issue number.

Without an issue number, anyone reading the code can see *what* the code does, not but not the *why*. Knowing the *why*
is crucial when reading code you don't understand - for understanding the code in itself, to decide if it can be
removed, make sure you can refactor it properly, and I'm sure other reasons. Especially as time goes by and the know-how
behind the code is forgotten, this becomes important.

If we follow this system, you can quickly get to the issue or PR of any code line doing git blame (or *Annotate* in
IntelliJ).

## Decision

<!--
    Short summary of the decision

    Example:
    We've decided to go for AlphaGov's system as explained [here](https://github.com/alphagov/govuk-aws/blob/24d1ea513e58ee938043d71d09815a51229067bf/docs/architecture/decisions/0001-record-architecture-decisions.md)
-->

Git commit messages found on the **master** branch should in **most cases** link to an issue, giving context to
the commit. The issue can be a `#`followed by Github issue number or `KM-` followed by a Trello card number.

Github issue number must be an issue or pull request (PR) identifier, so that anyone can put this message into the URL
`https://github.com/oslokommune/okctl/issues/<Github issue number>`. (Github automatically forwards to PR if issue is a
PR).

Examples:

```text
#123: Add userdocs for okctl create cluster
```

```text
#123: feat(ux): Add userdocs for okctl create cluster
```

```text
feat(ux): Add userdocs for okctl create cluster

#123
```

```text
KM-160: Add userdocs for okctl create cluster
```

```text
feat(ux): Add userdocs for okctl create cluster

KM-160
```

By *most cases*, we mean that there might be valid exceptions. We value pragmatism over too-rigid rules. The change
might be super-quick and small, where you know traceability doesn't add value. Use common sense. An example can be that
you just want to bump a version or change a typo.

Note that this RFC doesn't state anything about the commit messages in your feature branches. This RFC only cares about
the commits ending up on master. So you can adjust your workflow as you see fit, for instance squashing commits when
merging and adding on the PR's issue number to the commit message not before then.

## Consequences

<!--
    Short description of positive and negative consequences of the decision

    Example:
    The decisions we deem relevant will be documented in this format here in this folder.
-->

* We will get a git history that can be traced back to its background, giving more context.
* When browsing commit history on github.com, github automatically creates links to the github issue, so you'll be able
to click the "#123" part of the commit message to get directly to the issue.
* To find a Trello card given your number, put the card number into the search box.