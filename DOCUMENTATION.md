# School Management API Documentation

The School Management API Documentation, defining all functions and their purpose. As well as how to use the API, how to properly call functions and provide the necessary fields. 

### Contents

* [Getting Started](#getting-started)
* [General Functions](#general-functions)
* [Creating Accounts](#creating-accounts)
* [Log into Accounts](#log-into-accounts)
* [Updating Account information](#updating-accounts)

<br>

## Getting Started 

* ### Installation & Deployment
    The API is an all in one package and is simple and easy to deploy. Eliminating any long unforgiving          configurations you may face in other software.

    1. Clone the repo
        ```
        $ git clone https://github.com/SowinskiBraeden/school-management-api
        $ cd school-management-api
        ```
    
    2. Rename `.env.example` to `.env`
    
    3. Enter desired values into `.env`
        ```
        mongoURI='your mongo URI'
        dbo='school'
        secret='your 256 bit secret'
        
        # Suggested port for Production: 80
        # Suggested port for Development: 8000
        PORT='desired port'
        
        # This is to enable the system to send emails
        SYSTEM_EMAIL='your system email'
        SYSTEM_PASSWORD='your system email password'
        ```
    
    4. Run the system in your terminal
        ```
        $ go run main.go
        ```
    
    Or...
    
    5. Build and compile the system into an executable
        ```
        $ go build
        ```

* ### Initializing the System

    Upon running the system for the first time, you will be prompted to create an admin.
    In order to perform many actions with the API, an Admin account is required. You will be prompted to create your default admin account seen below...
    
    ![init](previews/init.png)
    
    Upon completing the initial setup and creating your default administrator for the system, it is ready to use and this box will appear below, displaying basic system details seen below...
    
    ![running](previews/running.png)
    
<br>

## General Functions

* ### API Status
    **Method:** `GET`
    ```
    <API_URL>/api/v1/status
    ```
    
    **Returns:**
    * Status 200: `OK`
    * JSON: 
        ```
        {
            "success": true,
            "message": "the API is active"
        }
        ```

<br>
  
## Creating Accounts

There are several account that can be registered into the system. As you may guess they are Administrators, Teachers and Students.

<br>

+ ### Creating Administrators
    You can create Administratos who will have permissions to perform majority of the actions in the API. It is common that there are more than one Administrator to help manage the system.

    **Method:** `POST`
    ```
    <API_URL>/admin/create
    ```
    
    **Required:**
    + Logged into an existing Admin account
    + JSON:
        ```
        {
            "firstname": "John",
            "lastname": "Doe",
            "dob": "01-01-1999",
            "email": "john_doe@example.com"
        }
        ```
    
    **Returns:**
    + Status 200: `OK`
    + JSON:
        ```
        {
            "success": true,
            "message": "successfully inserted admin"
        }
        ```
    
    <br>
    
+ ### Registering a Teacher
    Obviosuly a school management system will require teachers to manage and teach students. Teachers have an important role for the system.

    **Method:** `POST`
    ```
    <API_URL>/teacher/register
    ```
    
    **Required:**
    + Logged into an existing Admin account
    + JSON:
        ```
        {
            "firstname": "Homer",
            "middlename": "Jay",
            "lastname": "Simson",
            "dob": "12-05-1956",
            "email": "homerdog_simpson@example.com"
        }
        ```
    
    **Returns:**
    + Status 200: `OK`
    + JSON:
        ```
        {
            "success": true,
            "message": "successfully inserted teacher"
        }
        ```
        
+ ### Enrolling a Student
    What is a school withouth students? A professional day, but not the point. Students when their application to the school has been accepted by an admin can be enrolled into the school.

    **Method:** `POST`
    ```
    <API_URL>/student/enroll
    ```
    
    **Required:**
    + Logged into an existing Admin account
    + JSON:
        ```
        {
            "firstname": "Bart",
            "middlename": "JoJo",
            "lastname": "Simpson",
            "age": 10,
            "gradelevel": 4,
            "dob": "17-12-1979",
            "province": "...",
            "city": "Springfield",
            "address": "742 Evergreen Terrace",
            "postal": "..."
        }
        ```
        
    **Returns:**
    + Status 200: `OK`
    + JSON:
        ```
        {
            "success": true,
            "successfully inserted student"
        }
        ```
+ ### Additional Information After creating an account
    1. After successfully creating an account for another admin, teacher or student. A school email will be generated for them using their first and last name, each formated differently based on the type of account.

    2. All users are given a random ID used to sign into the system. Each ID is a random 6 digit number.
    
    3. Students are given a random PEN (Personal Education Number). A random 12 digit number.
    
    4. Accounts are given a defaul profile image. This can be updated in the future.
    
    5. All new users, once have verified their email, will recieve and email containing their ID.
    
    6. All new users, once have verified their email, will recieve and email containing their temporary password.
