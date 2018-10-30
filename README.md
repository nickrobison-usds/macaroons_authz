# Welcome to Buffalo!

Thank you for choosing Buffalo for your web development needs.

## Database Setup

It looks like you chose to set up your application using a postgres database! Fantastic!

The first thing you need to do is open up the "database.yml" file and edit it to use the correct usernames, passwords, hosts, etc... that are appropriate for your environment.

You will also need to make sure that **you** start/install the database of your choice. Buffalo **won't** install and start postgres for you.

### Create Your Databases

Ok, so you've edited the "database.yml" file and started postgres, now Buffalo can create the databases in that file for you:

	$ buffalo db create -a

## Starting the Application

Buffalo ships with a command that will watch your application and automatically rebuild the Go binary and any assets for you. To do that run the "buffalo dev" command:

	$ buffalo dev

If you point your browser to [http://127.0.0.1:3000](http://127.0.0.1:3000) you should see a "Welcome to Buffalo!" page.

**Congratulations!** You now have your Buffalo application up and running.

## What Next?

We recommend you heading over to [http://gobuffalo.io](http://gobuffalo.io) and reviewing all of the great documentation there.

Good luck!

[Powered by Buffalo](http://gobuffalo.io)

## Dependencies

### Gothic (for dev only)

We need gothic for user authentication, which comes with some nice plugings for buffalo.

If you need to make any changes, you can add the plugin.

```bash
go get -u github.com/gobuffalo/buffalo-goth
```

### Login.gov

This project requires a working, local installation of the Login.gov service.

You can clone the repo and build the docker images, like so:

```bash
git clone https://github.com/18F/identity-idp.git
cd identity-idp
bin/setup --docker
```

From there, you can issue the standard `docker-compose [up|down]` commands to get things running.
