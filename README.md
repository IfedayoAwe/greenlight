# greenlight
A fully featured JSON API for creating, retrieving and managing information about movies and users of the API.

## Description
A fast, secure, efficient, scallable and maintainable API server using Golang, a PostgreSQL database with migration files to make writing SQL migration easier and a Makefile to make running necessary development, code audit and build commands easier while also using the git hash to version build binaries. The idea behind this project is using dependency injection to implement efficient functions, using middleware and chaining handlers, taking into consideration best practices and writing efficient unit and integration tests, scallable project structure, vendoring third party packages, user authentication, securing the server and uses a PostgreSQL database but structured in a way that makes integration with other databases very easy. The project is designed as intended for usage in a production environment.

## Features
* Healthckeck: Shows enviroment, availablility status and application version
* User Registration: Automatic Movie Read Permissions and Default Profile Picture
* User Account Activation (Sends Token To Email)
* Create User Authentication/Login Token
* User Change Password
* Create User Reset Password Token (Sends Token To Email)
* User Reset Password
* Update User Details: Name, Email
* User Update Profile Picture
* Get User Details: Name, Email, Profile Picture Image-Path.
* Serving User Profile Picture
* User Logout
* Delete User Account
* List All Movies (Authenticated Users)
* Get A Specific Movie With It's ID (Authenticated Users)
* Add Movie Write Permissions For a User by an Admin
* Create A New Movie By Users With Movie Write Permissions
* Update A Movie 
* Delete A Movie
* Search For Movies Using Specific Query Parameters
* Dynamic Sorting For Movies Returned From The Database
* Dynamic Pagination For Movies Data
* Returning Movies Metadate (Current Page, Page Size, Total Pages, Total Records) with Movie Object 
* Request and Error Logging Using A Custom JSON Logger
* Unit and Integration tests
* Creating a Request Rate Limiter using users IP (Burst:4, r/s:2)
* Recover Panic
* Graceful Shutdown Of Application
* Configurable Request Origin Using Commandline Flags
* Displaying Application Metrics

## Installation
Clone the repo, set up a PostgreSQL database and execute the SQL migration files to create the necessary tables needed for the application to run and run **make run**

The following targets are defined in the Makefile:
```
- help                 lists the targets and their usage
- run                  run the cmd/api application
- startdb              runs a docker postgreSQL container
- createdb             creates a dabatase with user postgres
- dropdb               drops a database with user postgres
- migrateup            apply all up database migrations 
- migratedown          apply all down database migrations
- docker/compose/up    start containers in greenlight.yaml file
- docker/compose/down  stop and remove all running containers in greenlight.yaml file
- migration name=$1    create a new database migration file
- audit                tidy and vendor dependencies and format, vet and test all code
- vendor               tidy and vendor dependencies
- build/api            build the cmd/api application
- tests                runs test code coverage
```

## API Documentation
| Method |   URL Pattern              |  Action                                           |
|--------|----------------------------|---------------------------------------------------|
| GET    | /v1/healthcheck            | Show application health and version information   |
| GET    | /v1/movies                 | Show the details of all movies                    |
| POST   | /v1/movies                 | Create a new movie                                |                                                           
| GET    | /v1/movies/:id             | Show the details of a specific movie              |
| PATCH  | /v1/movies/:id             | Update the details of a specific movie            |
| DELETE | /v1/movies/:id             | Delete a specific movie                           |
| POST   | /v1/users                  | Register a new user                               |
| PUT    | /v1/users/activated        | Activate a specific user                          |
| PUT    | /v1/users/change-password  | Update the password of the request user           |
| PUT    | /v1/users/password         | Reset password of the request user                |
| Patch  | /v1/users/update-details   | Update the profile details of the request user    |
| PUT    | /v1/users/profile          | Update profile picture of the request user        |
| GET    | /v1/users/profile          | Get the profile details of the request user       |
| GET    | /profile/:filepath         | Serve Profile Picture                             |
| DELETE | /v1/users/logout           | Logout a user                                     |
| DELETE | /v1/users/delete           | Delete user account                               |
| POST   | /v1/tokens/activation      | Generate a new user activation token              |
| POST   | /v1/tokens/authentication  | Generate a new authentication token               |
| POST   | /v1/tokens/password-reset  | Generate a new password-reset token               |
| POST   | /v1/users/movie-permission | Give a user movie write permissions               |
| GET    | /debug/vars                | Display application metrics                       |

## Known Bugs
There are no known issues.

Test code coverage is around 62.6%.

## Contribution
Pull requests and new features suggestions are welcomed.