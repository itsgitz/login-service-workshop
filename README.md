# ORY KRATOS WORKSHOP

This repostory is contains Go source code for implementing Login System using ORY Kratos (User Management Provider), which consists of:

* Self-Service Browser Login System (Login using Web Browser Client)
* API Login System (Login using API Client)

# Usage

### A. Self-Service Browser
1. Run docker-compose for running ORY Kratos Service as User Management Service
```shell
$ make run
```

2. Access http://127.0.0.1:9080/login for landing on the login page, or http://127.0.0.1:9080/registeration for user registration (create a new user / sign up).

### B. API Client

Ory Kratos is not intended for API Client. According to their official website:

> API-based login and registration using this strategy will be addressed in a future release of ORY Kratos.
> https://www.ory.sh/kratos/docs/self-service/flows/user-login-user-registration/username-email-password#api-clients

I've tried to create API login system using Go HTTP Client, but this is still failed. ORY Kratos need cookies for sessions. CSRF Token is always invalid even all headers and cookies has copied into Go HTTP client that I've created, and login process is still failed.

# Contributor

Anggit M Ginanjar, Software Developer
