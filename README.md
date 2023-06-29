# go-webauthn

# Background

**Brief background on the problem statement**

In the field of user authentication, traditional password-based solutions have long been a security and user experience weak point. Passwords are vulnerable to a number of threats, including phishing attempts, reused passwords, and weak passwords. The Web Authentication (WebAuthn) standard, which offers a strong and password-free authentication mechanism based on public-key cryptography, was developed to overcome these problems.

**Why did I choose this problem?**

I chose this particular issue to solve the shortcomings and weaknesses of conventional password-based authentication techniques. Our initiative seeks to improve security, expedite authentication, and enhance the user experience. This decision supports greater security measures, promotes industry best practices, and increases consumer happiness and confidence in general.

**Why should others be interested in this problem?**

- **Improved Security:** WebAuthn increases security by removing weak passwords and weaknesses related to conventional authentication methods, lowering the risk of breaches and unauthorized access.
- **Improved User Experience:** Passwordless authentication makes the login process easier for users by eliminating the need to remember complicated passwords and reducing frustration with login-related difficulties, resulting in an improved user experience.
- **Privacy Protection:** WebAuthn uses public-key cryptography to prevent the transmission of user credentials to service providers, protecting user privacy and reducing the possibility of credential theft.
- **Compliance with Regulation:** Implementing WebAuthn demonstrates an organization's commitment to data privacy and compliance with legal obligations by aligning with data protection legislation like the GDPR.
- **Future-Proof Solution:** WebAuthn is a widely used standardized authentication protocol. It is supported by most major operating systems and web browsers and offers a future-proof solution that can adapt to new security threats and evolving authentication techniques.

## Problem Objectives

1. **Create a Golang backend server using the Mux router package and follow the RCM (Router, Controller, Model) model:**
   - Brief Explanation: The goal is to build a reliable and efficient backend server that can handle incoming requests, properly route them, and interact with the necessary models and controllers to perform the required actions.

2. **Implementation of User Models for Storing User Data in a No-SQL Database such as Couchbase on a Kubernetes Cluster:**
   - Brief Explanation: This task involves defining the necessary data structures and logic to store and retrieve user information in a No-SQL database. It ensures efficient data storage and retrieval for user-related operations.

3. **Construction of Credential Model to Securely Store Generated Credentials for each user in the Couchbase Server:**
   - Brief Explanation: This objective involves creating a model capable of securely storing and managing credentials, such as cryptographic keys or tokens, associated with each user. This ensures proper authentication and authorization during the login process.

4. **Creation of Session Model to Store Session Data for Active Sessions with Expiration Time:**
   - Brief Explanation: This objective focuses on creating a session model that enables the management and storage of session-specific data, such as user session tokens or session identifiers. It also includes handling session expiration to maintain security and efficiency.

5. **Development of APIs for User Registration Functionality:**
   - Brief Explanation: This task involves defining the necessary endpoints and implementing request/response handling mechanisms to facilitate user registration. It allows users to create new accounts by providing the required information, which is securely stored in the database.

6. **Creation of Minimalistic Frontend Interface for Handling Responses and Sending Registration Requests:**
   - Brief Explanation: This goal focuses on creating an intuitive interface where users can input their registration details, receive feedback on the registration process, and communicate with the backend server for registration-related tasks.

7. **Implementation of APIs for User Login and Validation of User Credentials:**
   - Brief Explanation: This task involves creating endpoints and implementing logic to authenticate user credentials, such as password verification or cryptographic token validation. The APIs generate appropriate responses based on the validation results.

8. **Enhance the frontend interface to handle login responses from the backend server and enable users to send login requests:**
   - Brief Explanation: This objective focuses on modifying the frontend interface to display login-related responses received from the server and allowing users to initiate login requests.

9. **Implement the necessary functions within the Couchbase User model to facilitate interaction with the database for CRUD operations (Create, Read, Update, Delete):**
   - Brief Explanation: This task involves writing functions within the User model that encapsulate the logic for interacting with the Couchbase database. These functions enable operations such as creating new user entries, retrieving user data, updating user information, and deleting user records.

10. **Integrate the defined routes, database interactions, and frontend components:**
    - Brief Explanation: This task involves connecting the implemented API routes with the corresponding database operations and ensuring the frontend interfaces appropriately interact with the backend server. It aims to establish a cohesive system where data flows smoothly between the frontend, backend, and database.

11. **Testing:**
    - Brief Explanation: This task involves designing and executing test cases to verify the correctness and robustness of the backend server, API endpoints, database interactions, and frontend functionality. It includes unit testing, integration testing, and potentially conducting security assessments to identify and address any vulnerabilities or bugs in the system.

## System Overview
<div  align="center">
  <img width="50%" src="https://github.com/gurpreet-legend/go-webauthn/assets/75157493/56904cf8-1d8c-40d6-a2a6-dbf4019edb6d"/>
  <img width="50%" src="https://github.com/gurpreet-legend/go-webauthn/assets/75157493/f3eb46ed-4b1d-4837-9239-3be0d271a49c"/>
</div>

## How to run
* Run command `cd client; npm run dev`, you could access the client at `localhost:5173`
* Follow these steps to set-up Couchbase server on Docker Desktop: [Couchbase Docker Desktop setup](https://docs.couchbase.com/server/current/install/getting-started-docker.html)
* Run the Couchbase server, you could access the server at `localhost:8091`
* At root of the project run `cd server; go run cmd/main.go`
