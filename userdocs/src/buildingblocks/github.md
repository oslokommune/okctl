`okctl` relies on services in AWS and Github to provide its functionality. In the following sections we describe some of the core services we use from Github.

## Github components

[Github](https://github.com/) is at its core a [Version Control System](https://en.wikipedia.org/wiki/Version_control) (VCS). Since its inception Github has grown to provide a wide array of functionality to interact with a product's source code. In Oslo kommune, we have decided to use Github as our primary source control system. It is therefore natural to make use of the functionality Github has on offer for implementing [Continuous Integration](https://en.wikipedia.org/wiki/Continuous_integration) (CI), [Continuous Deployment](https://en.wikipedia.org/wiki/Continuous_deployment) (CD), Authentication and Authorisation, within a product.

- [Actions](#github-actions-gha) for CI/CD
- [Organisations and Teams](#github-organisations-and-teams) for authorisation
- [Oauth Apps](#github-oauth-application-oauth) for authentication

### Github Actions (GHA)

[Github Actions](https://docs.github.com/en/actions) make it easy to automate, customize, and execute software development workflows from your repository. It is possible to discover, create, and share actions to perform any job one would like, including CI/CD, and combine actions in a completely customized workflow.

### Github Organisations and Teams

[Github Organisations](https://docs.github.com/en/github/setting-up-and-managing-organizations-and-teams/about-organizations) are shared accounts where businesses and open-source projects can collaborate across many projects at once. Owners and administrators can manage member access to the organization's data and projects with sophisticated security and administrative features.

[Github Teams](https://docs.github.com/en/github/setting-up-and-managing-organizations-and-teams/about-teams) allows one to organize members into teams that reflect the company or group's structure with cascading access permissions and mentions.

### Github Oauth Apps (Oauth)

[Github Oauth Apps](https://docs.github.com/en/developers/apps/building-oauth-apps) allow you to create and register an [OAuth App](https://oauth.net/2/) under your personal account or under any organization you have administrative access to. This Oauth app can then ensure that only the organisation or teams that you want to, have access to a UI or other protected service.
