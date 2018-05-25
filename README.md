# Mongo API Example
Inspired by [Restful API in Go with Mongodb](http://www.blog.labouardy.com/build-restful-api-in-go-and-mongodb/).

## Packages
* `github.com/BurntSushit/toml`
  * Reader for TOML file (may remove later)
* `github.com/globalsign/mgo`
  * Community supported MongoDB driver
* `github.com/goolge/uuid`
  * Generates UUIDs
* `github.com/julienschmidt/httprouter`
  * The HTTP router to create endpoints

## Version Control
Using [dep](https://github.com/golang/dep) for version control ([vgo](https://github.com/golang/go/wiki/vgo) has since been chosen for official version control).

## Cloud Ready
As part of this Go example, I wanted to see what the effort would be to deploy to Cloud Foundry, specifically [Pivotal Cloud Foundry](https://run.pivotal.io/).

### Manifest YAML
The manifest file is required for pushing the application to Cloud Foundry.

In the `manifest.yml`, it specifies the application's name, the buildpack (Cloud Foundry can figure out on its own), 
the memory (I am still testing the least amount of memory required), instances, environment variables, and any services to bind to.

#### Mongodb Service
This Go Mongodb API example requires a Mongo database to connect to. When deployed to Cloud Foundry, I 
am having the app bind to [mlab](https://docs.run.pivotal.io/marketplace/services/mlab.html). When the application 
binds to __mlab__, the application acquires connection information in an environment variable. I parse this 
information and use it to connect to the database.

### Pushing
To push to Cloud Foundry, change directories to where the `manifest.yml` file is located and run 
`cf push` (ensure you are logged into Cloud Foundry thru the __cli__ and have the __cli__ installed).

Cloud Foundry will build the application and run it.

## Local running
Simply run the `main()` function. 

Ensure to have a local instance of Mongodb running. The `config.toml` file determines the application to connect to the local DB 
(as I clean up the code, I will likey remove this to be 12 Factor compliant).