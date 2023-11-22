# Qlik Test Users

Backend service that can generate Qlik Sense tickets for specific test users

## UNDER DEVELOPMENT

The app is still being developed and issues might occur.

## Description

This idea had been on my mind for a long time now. On almost any Qlik Sense project there was an occasions that such service was needed. But always was pushed back because an workaround was found or was decided to invest the time for it on something else.

The basic idea is to have a list of test users in Qlik Sense and the service can generate web tickets for them. Once the ticket is available the developer/tester can use the ticket to login to Qlik Sense as the test user and "see" throw the eyes of that user. This service is not meant to be used in production environment and is to help the developers/testers to verify the work in the lower environments (pre-prod, UAT, DEV, Test etc)

The generated test users:

- can be added to section access with specific permissions (which can be altered)
- test the security rules behavior without constantly asking the real users to "refresh again?" in order to see the security rule implication

## Pre-Setup

### Executable

First we'll need the [latest release](https://github.com/informatiqal/qlik-test-users-tickets/releases)

### SSL Certificates

The app listen on `https` and we need certificates for that. If you dont have any at the moment the app itself can generate self-signed ones:

`.\qlik-test-users-tickets.exe --generateCert`

The command above will create `cert.pem` and `cert_key.pem` in the current folder

### Users

The app can generate set of test users, using the Repository API, in Qlik Sense. You can create as many as you want.

`.\qlik-test-users-tickets.exe --host <sense-host> --users "user1;user2;user3" --certPath c:\path\to\qlik\sense\certificates`

- `--host` - Qlik Sense host name. The command using Qlik Repository API so port `4242` should be accessible
- `--users` - semi-colon list of user ids. User ids can be anything as long as they do not repeat in the same user directory
- `--userDirectory` - we have to provide user directory to create the users under it. The app requires **suffix** and will create the users under `TESTING_<provided suffix>` user directory
- `--certPath` - path to a folder where Qlik Sense certificates can be located

### Config

Once we have the SSL certificates and the users created then the last thing, that is left, is to configure the app before installing the service and starting it.

The release archive will contain an example config - `config_example.toml`

```toml
[server]
port = 8081 # on which port the app will server the UI
httpsCertificatePath = "c:/path/to/https/pem/certificates" # full path to the SSL certificates

[qlik]
host = "machine-name" # Qlik Sense host name
certificatesPath = "c:/path/to/certificates" # full path to Qlik Sense certificates
userId = "sa_api" # user id, under which the app will communicate with Qlik
userDirectory = "INTERNAL"
domainName = "my-qlik.com" # used to generate the QMC and Hub links once the ticket is generated

[qlik.ports]
repository = 4242 # repository api port. default is 4242
proxy = 4243 # proxy api port. default is 4243

```

## Installation

At this point we are ready to install the service:

`.\qlik-test-users-tickets.exe --install`

The command will install the app as a Windows service. Once the service is installed it has to be started.

## Logging

The service produces two log files in the folder where the `qlik-test-users-tickets.exe` is located:

- app.log - general app log - server start, stop events. Any errors that have occurred. General audit - which tickets were generated
- http.log - the raw level http logs. Which endpoints were called, when, from which IP etc.

## Exposed endpoints

- `GET` `/healthcheck` - returns status `200`
- `POST` `/ticket` - this is the main endpoint. Its used to generate ticket for the provided user id:

  Request body:

  ```json
  {
      "userId":"1111111-2222-3333-4444-555555555555"
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

- `GET` `/virtualproxies` - used by the UI to return list with all possible virtual proxies
- `GET` `/users` - used by the UI to return list with the all available test users

## UI

The UI itself is very simple:

![UI](/assets/ui_image.png)

- choose test user from the list
- choose virtual proxy from the list
- press "Generate ticket" button
- get links to QMC and Hub with the ticket applied

## Future

I would like few (major) things to be implemented in the future releases:

- manage multiple Qlik Sense environments - at the moment the service can generate tickets for only one QS cluster and if there is a need for more clusters then the service have to be installed multiple times. Which is ok but I would like for the app be able to handle multiple clusters from one instance.
- manage multiple proxy services - some QS installations have multiple proxy services enabled and used and the virtual proxies are can be linked to these different proxies. This scenario require change to the config, the UI and a bit of backend
- attributes - the ticket generation endpoint accepts an optional property `attributes`. It will be useful (esp when talking about testing) to be able to provide additional attributes that will be associated the generated ticket. An example of attributes can be AD groups. Usually these attributes are set internally when Qlik generates the users session. But since we are using test users, that do not exists in any AD, the way to set these can be done via session/ticket attributes
