# Backend Server Written in Go

This is my prototype of REST API which is fulfilled in Go. It responses on messages from varies frontend web applications or mobile devices by using HTTP protocol and certaion commands which can be obtained in file `Tech Info.docx`. When resource is requested to this server and your client device has its necessary privileges, then you're up to get, modify, delete or create new entries in database resources.

To may have access to manage resources firstly sign in on system. Server will check your profile and will grant you access if your profile is activated. You can register more profiles but always awaits when superuser will activate your profile.

You can also use stream watch test function. The list of newly added links to `stream_config.json` file will be available on stream watch test page.

Stream recording is provided by sub-application named `auto_video_concat`. It concatenates video files into united files, and can be downloaded by this server via certain URL-request.

The stream watch test page is available by this address:

`127.0.0.1:8000/stream`

# Full List of Available URL Requests

* POST   /auth/signin
* POST   /auth/change_password
* POST   /auth/signout
* POST   /user
* GET    /user/:id
* GET    /user/all
* POST   /user/change_password
* PATCH  /user/:id
* DELETE /user/:id
* POST   /group
* GET    /group/:id
* GET    /group/all
* PATCH  /group/:id
* DELETE /group/:id
* POST   /group/user/:id
* GET    /group/user/:id
* DELETE /group/user/:id
* GET    /perm/all
* POST   /perm/user/:id
* GET    /perm/user/:id
* DELETE /perm/user/:id
* POST   /perm/group/:id
* GET    /perm/group/:id
* DELETE /perm/group/:id
* POST   /info
* GET    /info/:id
* GET    /info/all
* PATCH  /info/:id
* DELETE /info/:id
* POST   /video
* GET    /video/:id
* GET    /video/all
* PATCH  /video/:id
* DELETE /video/:id
* GET    /stream/get/:id
* GET    /stream/get/all

# How to use it

1. Create an `.env` file in directory `./configs/` and post variables from example `.env.example`.
2. Create database named "video_info" in your DBMS and create tables by executing
SQL query file named `up_database.sql`.
3. Build a binary with this command:
`go build -o ./build ./cmd/app`
4. You have to install or update several AV-libraries.
On Xubuntu 20.04 or higher post it to install/update all libraries:
`sudo apt-get install libavformat-dev`
`sudo apt-get install libavresample-dev`
`sudo apt-get install libavcodec-dev`
On Debian 11.3 or higher - install/update only two of those:
`apt install libavformat-dev`
`apt install libavresample-dev`
5. Make sure you have nginx installed by executing command:
`systemctl status nginx`
If it does not installed you have to install it by this command:
`apt install libpq-dev postgresql postgresql-contrib nginx curl`
