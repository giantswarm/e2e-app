[![CircleCI](https://circleci.com/gh/giantswarm/e2e-app.svg?&style=shield&circle-token=4a9231c3fd0adf22193a186149ccaeb3a72188d4)](https://circleci.com/gh/giantswarm/e2e-app)
[![Docker Repository on Quay](https://quay.io/repository/giantswarm/e2e-app/status "Docker Repository on Quay")](https://quay.io/repository/giantswarm/e2e-app)

# e2e-app

Application run for integration tests within a e2e test environment on Kubernetes.

The web application provides the following endpoints:

### `/`

Returns general info on the app, including the version, as a JSON object.

### `/delay/1`, `/delay/5`, `/delay/10`

Responds after 1, 5 or 10 seconds with a string `Hello World`.
