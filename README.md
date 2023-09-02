# BalkanID Bookshop - README

This project is a robust online bookstore designed for evaluation purposes for the BalkanID pre-interview round. It encompasses a wide range of features, including secure user authentication and authorization, account management, protection against security vulnerabilities, effective system logging, and a user-friendly interface for browsing and purchasing books. 

Built with a Go backend and a Next.js frontend, the project is containerized using Docker for easy setup. It includes a recommendation system to suggest books to users, ensuring a dynamic and personalized experience.

With thorough testing, deployment instructions, and comprehensive documentation, this project showcases coding proficiency, security awareness, and best practices in software development.

## Table of Contents

1. [Introduction](#introduction)
2. [Requirements](#requirements)
3. [Getting Started](#getting-started)
   - [Prerequisites](#prerequisites)
   - [Installation](#installation)
4. [Explanation](#explanation)
5. [Usage](#usage)
   - [User Guide](#user-guide)
6. [Testing](#testing)
    - [Unit Testing](#unit-testing)
    - [Workflow Testing](#workflow-testing)

## Introduction

This project is an online book store built with a Go backend and a Next.js frontend. It includes features for user authentication, book search and filtering, shopping cart management, account deactivation, and more. Additionally, it incorporates a recommendation system for suggesting books to users.

## Requirements

To run this project, you need the following:
- Just Docker 😅


## Getting Started

### Prerequisites

Before you start, ensure you have the required software installed on your system:

- Install Docker: [Docker Installation Guide](https://docs.docker.com/get-docker/)

### Installation

1. Clone the repository to your local machine:

   ```bash
   git clone https://github.com/BalkanID-University/vit-2025-summer-engineering-internship-task-ReuelNixon.git
   ```

2. Change into the project directory:

   ```bash
   cd bookshop
   ```

3. Build and start the application using Docker Compose:

   ```bash
   docker-compose build
   docker-compose up
   ```

   This command will set up the backend, frontend, and PostgreSQL database.

4. Access the application:
   - Frontend: http://localhost:8080
   - Backend: http://localhost:3000


## Explanation

1. **Backend**: The backend is built with Go and uses the Go Fiber framework. It includes a PostgreSQL database for storing user and book data and it is containerized using GORM. I chose Go Fiber because it has a lot of support online which can't be said for some other frameworks. 

2. **Frontend**: The frontend is built with Next.js and uses the Shadcn-UI library which is built on top of radix. I chose Next.js because it is a React framework that makes react development easier and faster. I chose Shadcn-UI because it is a UI library that is built on top of radix and it is very easy to use. 
3. **Docker**: While we build the docker image using docker-compose, it sets up the backend, frontend, and PostgreSQL database seperately. It pre-populates the database with some books. It also downloads the dependencies for the backend and frontend.

4. **Recommendation System**: Initially I thought of using python using a seperate docker container and exposing an API endpoint for getting recommendations. When I did little digging on recommendation algorithms, I got to know that user based colaborative filtering is the best approach for this use case. But before using this approach, we need to have a lot of data about the user to recommend books. So I decided to use a simpler approach. That is content based recommendation.  I built a ML model using the data of nearly 11000 books. I used the K-Nearest Neighbors algorithm to find the books that are similar to the book that the user is currently viewing. I used the book title, author to find the similarity between the books. Then I provided the user with the recommendations. But using a seperate docker container for python just to fetch similar books seemed to be a overkill to me. So I just dtored the data in a json file and used it in the go backend directly.

## Usage

### User Guide
The bookshop has an extensive collection of books, including both fiction and non-fiction. Some of the routes in the application are protected, so you will need to register for an account to access them. There are no pre-existing user accounts, so you will need to register for an account to access the bookshop.
- Non-protected activities:
    - Book list
    - Book search 
    - Book details
    - Book reviews
    - Book recommendations
- Protected activities:
    - Shopping cart
    - Checkout
    - Purchased books



## Testing
Testing is a critical aspect of the Online Book Store project, ensuring that the application functions correctly, reliably, and securely. The testing process includes both unit testing for individual components and workflow testing to evaluate the entire application's functionality.

### Unit Testing

The project includes a suite of unit tests that assess the functionality of individual components and functions within the backend. These tests cover various scenarios to validate the correctness of the code and identify potential bugs or issues.

### Workflow Testing

Integrated with GitHub Actions for continuous integration (CI). The workflow tests automatically run when code changes are pushed to the repository. These tests mimic the real-world deployment process, ensuring that the application behaves as expected in a production environment. 
