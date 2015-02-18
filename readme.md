# Glue

Lightweight language to easily handle communication between REST APIs

## Goals
 - Easy notation in a single file that outputs a binary. The binary broadcasts endpoints and returns/transfers the desired data. 
 - 

## Example:

Let's take the Google geocode API endpoint. We'll create a service that given a set of coordinates, returns a google maps image of that location. Let's take this endpoint: 
```
https://maps.googleapis.com/maps/api/geocode/json?address=1600+Amphitheatre+Parkway,+Mountain+View,+CA
```
From that endpoint we'll get a response like this:
```nginx
endpoint {

}
endpoint {

}
request {
    location: "/google.com"
}
```
```json
{
   "results" : [
      {
         "address_components" : [
            ...
         ],
         "formatted_address" : "1600 Amphitheatre Parkway, Mountain View, CA 94043, USA",
         "geometry" : {
            "location" : {
               "lat" : 37.4224879,
               "lng" : -122.08422
            },
            ...
         },
         "types" : [ "street_address" ]
      }
   ],
   "status" : "OK"
}
```

We want to: 
 - get a request
 - take the requests parameters
 - pass them to the google api
 - parse the response
 - pass it to the google image api string
 - return that as a response

First let's set up the request endpoint:

```nginx
server {
    location {
        endpoint /lat-lng-img;
        method get;
        request {
            method get;
            endpoint "https://maps.googleapis.com/maps/api/geocode/json";
            params {
                address params.address;
            }
        }
        return "https://maps.googleapis.com/maps/api/staticmap?center=" + 
            request.results[0].geometry.location.lat + 
            "," + 
            response.results[0].geometry.location.lng + 
            "&zoom=11&size=200x200";
    }
}
```

```json
"server": {
    "location": {
        "endpoint":"/lat-lng-img",
        "method":"get",
        "request " : [
            {
                "method" : "get",
                "url" : "https://maps.googleapis.com/maps/api/geocode/json",
                "params" : {
                    "address":"params.address"            
                }
            }
        ],
        "return":"'https://maps.googleapis.com/maps/api/staticmap?center=' + 
            request.results[0].geometry.location.lat + 
            ',' + 
            response.results[0].geometry.location.lng + 
            '&zoom=11&size=200x200'"
    }
}
```


```json
{
    "servers": [
        {
            "location": "/lat-lng-img"
        }
    ]
}
```