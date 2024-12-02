Gator is a multi-user CLI application. There's no server (other than the database), so it's only intended for local use.

Go and Postgres database will need to be installed

Create a config file in your home directory, ~/.gatorconfig.json, with the following content:
{
  "db_url": "postgres://example"
}

gator login - sets the current user in the config
gator register - adds a new user to the database
gator users - lists all the users in the database