# School Management API [![GoLang](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://go.dev/)

The backend API for a school management system accessible by staff and students.

***Explore the docs »***
* [Getting Started](#getting-started)
* [Usage/Documenation](#usage)
* [License](#license)
* [Contact](#contact)

### Built With

* [Golang](https://go.dev/)
* [GoFiber](https://gofiber.io/)
* [GoMongo Driver](https://docs.mongodb.com/drivers/go/current/)
* [golang-jwt](https://github.com/golang-jwt/jwt)

## About The Project

Originally started as a personal project to improve my overall Go knowledge and skills, driven by the complaints
of my current school management system for students to view their grades, credits, course selections etc. I took
this on as a full project to work on for my Software Engineering class and is now a staple of my work. I plan to
work on this till completion, starting with this API, and then the front end, learning new tricks and languages
along the way.

Developing the API has has me face many challenges. With me developing my techniques in Go, and increasing my 
ability to work with Go frameworks. Understanding the workflow of packages and the Go ecosystem.

I started by developing authentication; a solid footing to the project, without the authentication of users, 
the rest of the project would fail to work. As the purpose of this API handles sensitive data on students.

Once the authentication was complete, I took to working on majority of the user data handling. Such as 
setting, updating and or deleting user attributes. Building a few more key authentication features only 
allowing the proper authenticated user, an admin, perfom many of the actions.

Once this was complete I started work on conceptualizing the next major step of the API. Course management,
including course selection, schedule generation, etc. This can be viewed in this [request](https://github.com/SowinskiBraeden/school-management-api/issues/5).
To start off on this new set of features, I decided to tackle the largest and most complex problem, 
schedule generation. This algorithm generates a master timetable as well as updates the student schedule. 
All the while keeping track of any errors for admins to handle personally. This can be found in `testing`,
The [schedule generation](/test/scheduleGenerator) file is of the latest version, version 3. Though this has 
been moved to its own [repository](https://github.com/SowinskiBraeden/ScheduleGeneratorApp).

## Related Work

After the majority completion of vital api functionalities such as the ability to update information, create new
information, generate information, delete information, authentication etc. I have created a new repository 
for the next frontier of this project. I am now concurrently working on the [School Management Website](https://github.com/SowinskiBraeden/school-management-vue).

The School Management Website, is going to be written using the [Vue.js](https://vuejs.org/) framework.

## Getting Started

The systems is easy to start and your local machine. Eliminating any long unforgiving configurations you may face in other software.

### Installation

1. Clone the repo
```
    $ git clone https://github.com/SowinskiBraeden/school-management-api.git
    $ cd school-management-api
```

<br>

2. Rename `.env.example` to `.env`

<br>

3. Enter desired values into `.env`
```
    mongoURI='your mongo URI'
    dbo=school
    secret='your 256 bit secret'
    PORT='your desired port'
    SYSTEM_EMAIL='your system email'
    SYSTEM_PASSWORD='your system email password'
```

<br>

4. Run the system in your console
```
  go run main.go
```

* The first time the API is run, you will be promted to enter in the systems defualt administrator account detials.
 This is required, in order to use majority of the API you'll need to be in an authenticated admin account.

<br>

5. Build and compile the system into an executable
```
  go build
```
<br>


## Usage

The usage for this API is in the name, school management, it's not an easy task in the areas I live near. Overcrowded schools are far too common and the need for a better school management system is there. The current systems issued by the government work but are poor and outdated.

To see a full list of features and how to use the system, read the [documentation](DOCUMENTATION.md)


<!-- LICENSE -->
## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.


<!-- CONTACT -->
## Contact

Braeden Sowinski - [@BraedenSowinski](https://twitter.com/BraedenSowinski) - sowinskibraeden@gmail.com - McDazzzled#5307 on Discord

Project Link: [https://github.com/SowinskiBraeden/school-management-api](https://github.com/SowinskiBraeden/school-management-api)
