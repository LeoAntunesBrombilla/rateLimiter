#!/bin/bash

# URL of your server
url="http://localhost:8080"

# Make a number of requests
for i in {1..20}
do
   curl $url
done