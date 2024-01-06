#!/bin/bash

# URL of your server
url="http://localhost:8080"

# Token value
token="10"

# Make a number of requests
for i in {1..20}
do
   curl -H "API_KEY: $token" $url
done