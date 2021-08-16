# Bigquery Importer for Slack

Go service for importing Slack workspace data into BigQuery

## Usage

To use th service, the following environment variables have to be set:

Variable Name | Description
---|---
LOGGER_SERVICENAME | Will add the ServiceContext to the log with the specified service name.
LOGGER_LEVEL| The minimum enabled logging level. Recommended: **debug**.
LOGGER_DEVELOPMENT | If the logger is set to development mode or not. Recommended: **false**.
SLACKCLIENT_APISECRET | The service requires that the API key for accessing the Slack workspace data is stored in a Secret Manager secret. This variable should be set to the full resource name of that secret.
BIGQUERYCLIENT_PROJECTID | The id of the project where the tables will be created.
JOB_DATASET | The name of the dataset where the tables will be created.
JOB_ORG | The organization the data belongs to.
JOB_APPENDIDSUFFIX | When this flag is true the job's id will be used as a suffix for the table name. This is useful for testing when multiple tables have to be created in quick succession. Recommended: **false**.

The Slack API Key is acquired by creating and installing a new Slack bot on the workspace that will have its data
exported. Instructions can be found [here][slack-api-key]. The key should be of the bot-token type and contain the
following scopes:

- channels:read
- groups:read
- usergroups:read
- users:read
- users:read.email
- files:read

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to add and update tests as appropriate.

Contributions should adhere to the [Conventional Commits][commits] specification.

## License

[MIT](https://choosealicense.com/licenses/mit/)

[slack-api-key]:https://api.slack.com/authentication/token-types#bot

[commits]:https://www.conventionalcommits.org/en/v1.0.0/
