# School Management API [![GoLang](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://go.dev/)

The backend API for a school management system accessible by staff and students.

<a href="https://github.com/SowinskiBraeden/school-management-api"><strong>Explore the docs »</strong></a>
<br>
<br>

## About The Project

Originally started as a personal project to improve my overall Go knowledge and skills, driven by the complaints
of my current school management system for students to view their grades, credits, course selections etc. I took
this on as a full project to work on for my Software Engineering class and is now a staple of my work. I plan to
work on this till completion, starting with this API, and then the front end, learning new tricks and languages
along the way.

## Future Plans

After the majority completion of vital api functionalities such as the ability to update information, create new
information, generate information, delete information, authentication etc. I plan on creating a new repository 
for a new frontier of this project. I plan on making the [School Management Website](https://github.com/SowinskiBraeden/school-management-website).

This is most likely going to be written in another programming language or framework I am unfamiliar with. The 
idea is to have a modern designed website capable of being used easily on any device or platform. I may use the
[React.js](https://reactjs.org/) or [Vue.js](https://vuejs.org/) frameworks to build the front end design.
<br>

### Built With

* [Golang](https://go.dev/)
* [GoFiber](https://gofiber.io/)
* [GoMongo Driver](https://docs.mongodb.com/drivers/go/current/)

## Getting Started

The systems is easy to start and your local machine. Eliminating any long unforgiving configurations you may face in other software.

### Installation

1. Clone the repo
```
    $ git clone https://github.com/SowinskiBraeden/school-management-api.git
    $ cd school-management-api
```

2. Rename `.env.example` to `.env` 
3. Enter desired values into `.env`
   ```
    mongoURI='your mongo URI'
    dbo=school
    secret='your secret'
    PORT='your desired port'
    SYSTEM_EMAIL='your system email'
    SYSTEM_PASSWORD='your system email password'
   ```
4. Run the system in your console
```
  go run main.go
```
5. Build and compile the system into an executable
```
  go build
```
<br>


## Usage

The usage for this API is in the name, school management, it's not an easy task in the areas I live near. Overcrowded schools are far too common and the need for a better school management system is there. The current systems issued by the government work but are poor and outdated.


<!-- CONTRIBUTING -->
## Contributing

Please refer to [CONTRIBUTING](CONTRIBUTING.md) for contributing to the project.



<!-- LICENSE -->
## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.


<!-- CONTACT -->
## Contact

Braeden Sowinski - [@BraedenSowinski](https://twitter.com/BraedenSowinski) - sowinskibraeden@gmail.com

Project Link: [https://github.com/SowinskiBraeden/school-management-api](https://github.com/SowinskiBraeden/school-management-api)
