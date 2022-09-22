## Summary

Note: We most likely forgot writing this post mortem properly, so this file is here as a placeholder to not
forget that we actually had an incident. We need to have this for proper insight into our production
incident.

TL;DR: 

Loki stopped working due to the Loki pod was missing an IAM permission "DeleteTable".

Details:

* Communication with other teams about the same issue:
  * https://oslokommune.slack.com/archives/CV9EGL9UG/p1661499598067769
  * https://oslokommune.slack.com/archives/CV9EGL9UG/p1662033302833479
