# CI/CD Standard Proposal: 0.0.0

A CI/CD pipeline design and configuration standard proposition defined using yaml format.
---

## The Problems

1. **Every business with software development teams tries to build a CI/CD, but no standard exists.**

2. **Most businesses don't have the resources, talents and velocity to build a top of the line CI/CD pipelines system.** 

3. **Vendor lock-in of the different available tools kills the velocity of ecosystem.**

## Proposed Solution

Following the Agile best practices, developing with iterations. In this case, pipelines. This is a proposition of a standerdized interface, plain and simple, inspired by the CRI (*Container Runtime Interface) and other similar standards, but in this case specific to CI/CD solutions. It also comes with best practices recommendations to facilitate the implementation.

Because it's a simple and widely used markup language, this proposition uses the yaml format in a Kubernetes inspired style to define the interface.

## Goals

- **Improve reusability of CI/CD configurations and integrations**
- **Improve feature development velocity**
- **Improve code maintainability**
- **Improve best practices**
- **Tools, frameworks and languages agnostic**

## Non-Goals

- **Restrict developers to a specific workflow**
- **Restrict features**
- **Cover unicorn use cases**

## CI/CD Interface

A CI/CD pipeline is composed of a multitude of services and its user usually defines a series of tasks calling those services. A task can therefore be represented as a set of parameters sent to a service which will produce a result.

Some parameters are required for a multitude of tasks, with the only exception of simple tasks. Tasks can be categorized by well known actionable concepts such **build**, **release** and **deploy**. Those concepts paired with the parameters will serve as the baseline of the proposed interface.

*Leaving out the plan, code, test, operate and monitor concepts for the next iteration.*

## Tasks General Rules

```yaml
# Tasks definition
  # Tasks represent code or programs to be executed with the required parameters
# Usage rules:
  # 1. Variables MUST NOT be made of other variables. This should be handled within the service integration
# Implementation rules:
  # 1. Tasks MUST call at least one service
  # 2. Tasks MUST produce a result
  # 3. Tasks results object has 3 sub-objects: 
  #       the result code(integer), the execution log artifact (url), the pipeline artifacts (url)

product_data:
  name: # Name of the product

tasks:
# Example build type task
- name: build_app_a # The task name, MUST be unique
  type: build # The task type
  after: [] # The tasks required to be done and successful before this one
  build:
    environment: my-docker-image # The environment in which to build the product
    command: docker build $REPOSITORY_DIR # The command to be executed to build the product 
                                          # If multiple commands are required, then it should be a scipt
  # Implementation rules: 
    # 1. The build MUST produce one or many build artifacts
    # 2. Only a build CAN produce a build artifact
    # 3. The build artifacts MUST have a unique build identifier per artifact, known as the build number


# Example release type task
- name: alpha_release_app_a
  type: release # The task type
  after: ['build_app_a'] # The tasks required to be done and successful before this one
  release:
    level: alpha # The release level of this task, possible values based on **semantic versioning
    type: MINOR # The impact of this release, possible values based on **semantic versioning
    metadata: # [Optional]
  # Implementation rules: 
    # 1. A Release task MUST create a release or pre-release version that will be associated
    #    with the corresponding build number
    # 2. A Release task MUST NOT produce new build artifacts. It MAY tag upstream build artifacts
    # 3. Versioning SHOULD implement **semantic versioning 2.0.0 standard

- name: deploy_app_a_to_prod
  type: deployment
  after: ['alpha_release_app_a']
  deployment:
    environment: staging # [Required] the target environment 
    release:
      task: alpha_release_app_a  # [Required] Reference a previous release task
  # Implementation rules: 
    # 1. The deployment task MUST only update a desired state database (example: git repository)
    #    for the target environment. It MAY trigger a synchronisation of the orchestrator with the database
    # 2. Build artifacts MUST be pulled by the orchestrator
    # 3. Deployment configurations SHOULD be stored and versioned within the deployed application's repository
    ## Note: In this context, an orchestrator is a service able to read the desired state database
    ##       and apply/deploy the desired state, example: ***ArgoCD. More on this topic in future article

```

*\* [Container Runtime Interface proposal](https://github.com/kubernetes/kubernetes/blob/release-1.5/docs/proposals/container-runtime-interface-v1.md)*
*\*\* [Semantic versioning](https://semver.org)*
*\*\*\*[ArgoCD](https://argo-cd.readthedocs.io/en/stable/)*
## Contribution

For questions, coments and suggestions open a GitHub issue or pull request.
## LICENSE
<a rel="license" href="http://creativecommons.org/licenses/by/4.0/"><img alt="Creative Commons License" style="border-width:0" src="https://i.creativecommons.org/l/by/4.0/80x15.png" /></a><br />This work is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by/4.0/">Creative Commons Attribution 4.0 International License</a>.


