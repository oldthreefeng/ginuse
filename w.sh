#!/bin/sh
data=$(date +%F-%T)
### blog is running by `hugo server -w`
### when push something to this git repo, this will be executed
### do something to deploy.
cd /app/Blog && git pull 

### append the exec into log file
echo "$data execute shell" >> /app/blog.log
