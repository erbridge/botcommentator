# gotwit

A library for creating twitter bots.


## Getting Started

Copy the files from `template/` to a new project.

Rename `secret.example.json` to `secrets.json` and fill in the tokens with your bot's API tokens.


## Running on Heroku

Install Godep:

```
$ go get github.com/tools/godep
```

Freeze the dependencies:

```
$ godep save
```

Add a `Procfile` containing the following:

```
worker: <project-name>
```

Commit everything:

```
$ git init
$ git add .
$ git commit -m "Initial commit."
```

Now you're ready to push to Heroku.

Create a Heroku instance:

```
$ heroku create -b https://github.com/heroku/heroku-buildpack-go.git <project-name>
```

Set up the environment:

```
$ heroku config:set \
    GOTWIT_CONSUMER_KEY=<consumer-key> \
    GOTWIT_CONSUMER_SECRET=<consumer-secret> \
    GOTWIT_ACCESS_TOKEN=<access-token> \
    GOTWIT_ACCESS_TOKEN_SECRET=<access-token-secret>
```

Now you're ready to go. Push it live:

```
$ git push heroku master
```
