#  GOINVEST-UP

A Cloud Native application that aims to manage stocks & predicts opening price of stocks. This repository has the Rest Api codebase of the application.
The API uses [chi](https://github.com/go-chi/chi) router; [gorm](https://github.com/jinzhu/gorm) ORM for Connection with the MYSQL database and 
The application uses [JWT](https://github.com/go-chi/jwtauth) authentication middleware for user authentication. The Application has been Containerized using Docker and the image of the application can be pulled from Dockerhub.
       
## Getting Started

These instructions will get you a copy of the api and running on your local machine for development and testing purposes.

### Installing

For just running the API you can simply :

1) Pull the Docker Image

```
docker pull arc1999/goinvest
```

2) set the environment variables in a .env file

```
CONTAINER_NAME=cloudserver
AWSDB= username:password@tcp(awsdbendpointhere
```


3) Run the Docker-compose file.

```
until finished
```

For Testing Purposes you might Clone the repo and repeat step 2 & step 3 of the above installation.

## Running the tests

The Test Cases of API are still in Progress.

## Deployment

The Api has been Containerized using Docker and can be deployed using any 

## Built With

* [Dropwizard](http://www.dropwizard.io/1.0.2/docs/) - The web framework used
* [Maven](https://maven.apache.org/) - Dependency Management
* [ROME](https://rometools.github.io/rome/) - Used to generate RSS Feeds

## Contributing

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/your/project/tags). 

## Authors

* **Billie Thompson** - *Initial work* - [PurpleBooth](https://github.com/PurpleBooth)

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Hat tip to anyone whose code was used
* Inspiration
* etc
