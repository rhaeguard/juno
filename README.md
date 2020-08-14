<p align="center">
  <img src="./.assets/juno.png" alt="juno logo" width="250px"/>
</p>

Juno is a very tiny and simple to use http file server

## Endpoints

| method    | endpoint                  | does                                                                      |
|:---       |:---                       | :---                                                                      |
|POST       |/v1/auth/login             | returns access token                                                      |
|POST       |/v1/auth/refresh_token     | refreshes the access token                                                |
|POST       |/v1/auth/logout            | invalidates the token                                                     |
|GET        |/v1/resources              | retrieves all the resources related to the application                    |
|POST       |/v1/resources/upload       | uploads the given file                                                    |
|GET        |/v1/resources/:id          | retrieves a single resource information or downloads that file            |
|DELETE     |/v1/resources/:id          | deletes all the information related to the resource with the given id     |

## Launching

Run the following command to launch the application

```shell script
docker-compose up
```