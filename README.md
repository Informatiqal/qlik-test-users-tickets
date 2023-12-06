# Qlik Test Users

Backend service that can generate Qlik Sense tickets for specific test users

## Description

This idea had been on my mind for a long time now. On almost any Qlik Sense project there was an occasions that such service was needed. But always was pushed back because an workaround was found or was decided to invest the time for it on something else.

The basic idea is to have a list of test users in Qlik Sense and the service can generate web tickets for them. Once the ticket is available the developer/tester can use the ticket to login to Qlik Sense as the test user and "see" through the eyes of that user.

### This service is not meant to be used in production environments. It is to help the developers/testers to verify the work in the lower environments (pre-prod, UAT, DEV, Test etc)

The generated test users:

- can be added to section access with specific permissions (which can be altered)
- test the security rules behavior without constantly asking the real users to "refresh again?" in order to see the security rule implication

## Pre-Setup

### Executable

First we'll need the [latest release](https://github.com/informatiqal/qlik-test-users-tickets/releases/latest)

### SSL Certificates

The app listen on `https` only and we need certificates for that. If you dont have any at the moment the app itself can generate self-signed ones:

`.\qs_test_users.exe gencert`

The command above will create `cert.pem` and `cert_key.pem` in the current folder

### Users

The app can generate set of test users, using the Repository API, in Qlik Sense. Using the data from `config.toml` its possible to create and many as needed at any time.

`.\qs_test_users.exe users create --suffix SOMETHING --users "user1;user2;user3..." --cluster "some cluster name here"`

- `--users` - semi-colon list of user ids. User ids can be anything as long as they do not repeat in the same user directory. Ideally wrap the users string into double-quotes
- `--suffix` - we have to provide user directory to create the users under it. The app requires **suffix** and will create the users under `TESTING_<provided suffix>` user directory
- `--cluster` - name of the cluster on which the users will be generated

### Config

Once we have the SSL certificates and the users created then the last thing, that is left, is to configure the app before installing the service and starting it.

The release archive will contain an example config - `config_example.toml`

```toml
[server]
port = 8081 # where to run the app. If not provided 8081 will be used
httpsCertificatePath = "c:/path/to/https/pem/certificates" # full path to the SSL certificates

# define at least one cluster
[qlik."some cluster name"]
certificatesPath = "c:/path/to/certificates" # full path to Qlik Sense certificates
repositoryHost = "machine-name" # host/machine name where the repo communication to be made
userId = "sa_api" # (optional) on behalf of which user the communication to be made. default is sa_api
userDirectory = "INTERNAL" # (optional) internal user userDirectory. default is INTERNAL
trustAllCerts = true # (optional) to ignore or not all ssl communication issues (like self-signed certificates). default is true

# optional
# mapping between the existing host/machine names and the pretty url
# when the ticket is generated apply this mapping for the return QMC and Hub urls
# if this section do not exists or no mapping is found
# then the node's machine name will be used
[qlik."some cluster name".domainMapping]
"machine-name1" = "my-qlik.com"
"another-machine-name" = "my-other-domain.com"

# define as many clusters as needed if the app will serve more than one cluster
# the UI have an option to choose on which cluster the ticket will be created
[qlik."another cluster"]
...
[qlik."another cluster".domainMapping]
...
```

## Installation

At this point we are ready to install the service:

`.\qs_test_users.exe service install`

The command will install the app as a Windows service. Once the service is installed it has to be started.

P.S. To uninstall the service simple run:

`.\qs_test_users.exe service uninstall`

> [!TIP]  
> The command should be executed with administrator privileges. Otherwise you will get "access denied" error

## Logging

The service produces two log files in the folder where the `qs_test_users.exe` is located:

- `app.log`- general app log - server start, stop events. Any errors that have occurred. General audit - which tickets were generated
- `http.log` - the raw level http logs. Which endpoints were called, when, from which IP etc.

## Exposed endpoints

- `GET` `/healthcheck` - returns status `200`
- `POST` `/ticket` - this is the main endpoint. Its used to generate ticket for the provided user id:

  Request body:

  ```json
  {
      "userId": "1111111-2222-3333-4444-555555555555"
      "proxyId": "6666666-7777-8888-9999-000000000000" // id of the proxy service. The ticket will be generated on the host that is behind that proxy service
      "cluster": "some cluster name" // name of the cluster on which the ticket to be generated. The names are defined in config.toml
      "virtualProxyPrefix": "something" // optional. If not provided then the service will generate the ticket for the default ("/") virtual proxy
  }
  ```

  Response:

  ```json
  {
    "userId": "1111111-2222-3333-4444-555555555555",
    "userDirectory": "TESTING_SOMETHING",
    "ticket": "1234567890123456",
    "virtualProxyPrefix": "something",
    "links": {
      "qmc": "https://my-sennse.com/something/qmc?qlikTicket=1234567890123456",
      "hub": "https://my-sennse.com/something/hub?qlikTicket=1234567890123456"
    }
  }
  ```

- `GET` `/proxies/:id` - used by the UI to return list with all active proxies and their virtual proxies on the specified cluster `:id`
- `GET` `/users/:id` - used by the UI to return list with the all available test users on the specified cluster `:id`
- `GET` `/clusters` - used by the UI to return list of the available/setup cluster names

## CLI commands

Apart from [users create](#users) and [gencert](#ssl-certificates) command the other useful command is `ticket create`

This command allows us to create ticket but from the terminal (although still thinks that through the UI is easier).

`.\qs_test_users.exe ticket create --cluster "some cluster name" --proxyId 6666666-7777-8888-9999-000000000000 --suffix "SOMETHING" --userId "some user id"`

Have a look at the command help for more info `.\qs_test_users.exe ticket create --help`

## UI

The UI itself is very simple:

![UI](/assets/ui_image.png)

- choose cluster (if only one cluster is set then it will be pre-selected)
- choose test user from the list
- choose proxy from the list (if only one proxy exists then it will be pre-selected)
- choose virtual proxy from the list (after proxy is selected)
- press "Generate ticket" button
- get links to QMC and Hub with the ticket applied

## Future

**Update (2023-12-06)**: all of the initial major features are now implemented. Apart from some code refactoring (backend and frontend) no major features are planned ... for now ;)

I would like few (major) things to be implemented in the future releases:

- **implemented** ~~manage multiple Qlik Sense environments - at the moment the service can generate tickets for only one QS cluster and if there is a need for more clusters then the service have to be installed multiple times. Which is ok but I would like for the app be able to handle multiple clusters from one instance.~~
- **implemented** ~~manage multiple proxy services - some QS installations have multiple proxy services enabled and used and the virtual proxies are can be linked to these different proxies. This scenario require change to the config, the UI and a bit of backend~~
- **implemented** ~~attributes - the ticket generation endpoint accepts an optional property `attributes`. It will be useful (esp when talking about testing) to be able to provide additional attributes that will be associated the generated ticket. An example of attributes can be AD groups. Usually these attributes are set internally when Qlik generates the users session. But since we are using test users, that do not exists in any AD, the way to set these can be done via session/ticket attributes~~
