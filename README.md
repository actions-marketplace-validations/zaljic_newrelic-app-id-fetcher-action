# NewRelic app id fetcher

[![Actions Status](https://github.com/zaljic/newrelic-app-id-fetcher-action/workflows/Build/badge.svg)](https://github.com/zaljic/newrelic-app-id-fetcher-action/actions)

This action fetches an app id of an existing service from the NewRelic applications API.

> NOTE: The app to fetch the ID of must already exist in the NewRelic account. The action does not create a new app.

## Usage

To use the action, the NewRelic API key must be provided as a secret in the repository.

### Example workflow

```yaml
- name: Fetch NewRelic app id
  id: newrelic-app-id
  uses: zaljic/newrelic-app-id-fetcher-action@v1
  with:
    newrelicApiKey: ${{ secrets.NEWRELIC_API_KEY }}
    newRelicRegion: EU
    appName: my-app-name
```

### Inputs

| Input                                             | Description                                        |
|------------------------------------------------------|-----------------------------------------------|
| `newrelicApiKey`  | The NewRelic API Key of the account the app id should be retrieved from    |
| `newrelicRegion` _(optional)_ | The region of the NewRelic account the app is monitored in. Defaults to  `US`   |
| `appName`  | The name of the app to fetch the ID of    |

### Outputs

| Output                                             | Description                                        |
|------------------------------------------------------|-----------------------------------------------|
| `appID`  | The ID of the app specified in `appName`    |

## Examples

> NOTE: The examples don't work yet. The app id is required to fetch the GUID of the app from the NewRelic API. The GUID action is not yet implemented.

### Using the optional input

This is how to use the optional input.

```yaml
with:
  newrelicRegion: EU
```

### Using outputs

You can use the output of this action to fetch an app GUID from the NewRelic API and use it as input for the [newrelic/deployment-marker-action](https://github.com/newrelic/deployment-marker-action)

```yaml
name: newrelic

on:
  workflow_dispatch:

jobs:
  newrelic:
    name: newrelic
    runs-on: ubuntu-latest
steps:
  - uses: actions/checkout@master
  - name: Fetch NewRelic app id
    id: newrelic-app-id
    uses: zaljic/newrelic-app-id-fetcher-action@v1
    with:
      newrelicApiKey: ${{ secrets.NEWRELIC_API_KEY }}
      newRelicRegion: EU
      appName: my-app-id

  # Use the output from the `newrelic-app-id` step to fetch the GUID of the app
  - name: Fetch NewRelic app GUID
    id: newrelic-app-guid
    uses: zaljic/newrelic-app-guid-fetcher-action@v1
    with:
      newrelicApiKey: ${{ secrets.NEWRELIC_API_KEY }}
      newRelicRegion: EU
      appID: ${{ steps.newrelic-app-id.outputs.appID }}

  # Use the output from the `newrelic-app-guid` as input for the newrelic deployment-marker-action
  - name: Create NewRelic deployment marker
    uses: newrelic/deployment-marker-action@v2
    with:
      apiKey: ${{ secrets.NEW_RELIC_API_KEY }}
      guid: ${{ steps.newrelic-app-guid.outputs.appGUID }}
      version: ${{ github.sha }}
      changelog: ${{ github.event.head_commit.message }}
      description: release
      region: EU
      user: ${{ github.actor }}
```