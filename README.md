# Simple Golang Demo

Entities:
- Account
- AccountBalance
- AccountTransaction

Features: 
- Create account
- Get account by ID
- Delete account by ID
- Update account
- Fund transfer

# How to run

- Run `docker-compose -f docker-compose-dependencies.yaml up -d`
- Go to `localhost:8888` to access PGAdmin
- Create database with name `mockva`
- Create file `.env`, please refer to `.env.example`
- Run project