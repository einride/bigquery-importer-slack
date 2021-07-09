# Bigquery Importer for Slack

Go service for importing Slack workspace data into BigQuery

## Usage

To use th service, the following environment variables have to be set:

- **LOGGER_SERVICENAME**: Will add the ServiceContext to the log with the specified service name.
- **LOGGER_LEVEL**: The minimum enabled logging level. Recommended: Debug
- **LOGGER_DEVELOPMENT**: If the logger is set to development mode or not. Recommended: false
- **SLACKCLIENT_APISECRET**: The service assumes that the Slack API key is stored in a GCP Secret Manager. This variable
  should therefore be set to the full resource path of that secret. The API Key is acquired by creating and installing a
  new Slack bot on the workspace that will have its data exported. Instructions can be found [here][slack-api-key]. the
  key sould be of the bot-token type and contain teh following scopes:
    - channels:read
    - groups:read
    - usergroups:read
    - users:read
    - users:read.email
- **BIGQUERYCLIENT_PROJECTID**: The project that contains the BigQuery Dataset
- **JOB_DATASET**: The name of the dataset that the data will be exported to.
- **JOB_ORG**: The organisation the data belongs to.

[slack-api-key]:https://api.slack.com/authentication/token-types#bot

## License
[MIT](https://choosealicense.com/licenses/mit/)
