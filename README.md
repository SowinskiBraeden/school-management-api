# School Management API [![GoLang](https://pkg.go.dev/badge/golang.org/x/example.svg)](https://go.dev/)

The backend API for a school management system accessible by staff and students.

***Explore the docs »***
* [Getting Started](#getting-started)
* [Usage/Documenation](#usage)
* [Contributing](#contributing)
* [License](#license)
* [Contact](#contact)
<br>
<br>

## About The Project

Originally started as a personal project to improve my overall Go knowledge and skills, driven by the complaints
of my current school management system for students to view their grades, credits, course selections etc. I took
this on as a full project to work on for my Software Engineering class and is now a staple of my work. I plan to
work on this till completion, starting with this API, and then the front end, learning new tricks and languages
along the way.

I highly recommend looking at the testing for [schedule generation](/test/scheduleGenerator), version 3.
Though this has been moved to its own [repository](https://github.com/SowinskiBraeden/schedule-generator).

## Related Work

After the majority completion of vital api functionalities such as the ability to update information, create new
information, generate information, delete information, authentication etc. I have created a new repository 
for the next frontier of this project. I am now concurrently working on the [School Management Website](https://github.com/SowinskiBraeden/school-management).

The School Management Website, is going to be written using the [Next.js](https://nextjs.org/) framework. Previously the plan  
was to use Vue.js, and I started work on the project, which I have now decided to discontinue and archive it [here](https://github.com/SowinskiBraeden/school-management-archived). I now plan to start from scratch and use Next.js as I believe it is more  
suited for this project and is more popular with its use of React.js over Vue.js. The new repository can be found [here](https://github.com/SowinskiBraeden/school-management).
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
