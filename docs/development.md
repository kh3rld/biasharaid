# Development Guide

This section provides instructions and guidelines for developers contributing to the `biasharaid` project. Follow the steps and guidelines outlined below to set up your development environment, adhere to the project's code style, and understand the contributing process.

## Setting Up

To set up the development environment for the `biasharaid` project, follow these steps:

### 1. Clone the Repository

Clone the repository from GitHub to your local machine.

```bash
git clone https://github.com/kh3rld/biasharaid.git
cd biasharaid
```

### 2. Install Dependancies 

Ensure you have Go installed on your machine. Install project dependancies using `go mod' commands.

```sh
go mod tidy
```

### 3.  Running the Server

To run the server locally, use the `go run` command.

```sh
go run main.go
```
The server will be running at `http://localhost:8080`.

### 4. Directory Structure

The project directory structure is as follows:
```
.
├── blockchain
│   └── blockchain.go
├── internals
│   └── routes
│       ├── handlers.go
│       └── routes.go
├── renders
│   ├── render.go
└── views
    ├── static
    │   ├── 404.css
    │   ├── 500.css
    │   ├── contact.css
    │   ├── dummy.css
    │   ├── index.css
    │   └── verify.css
    └── templates
        ├── 404.page.html
        ├── 500.page.html
        ├── add.page.html
        ├── base.layout.html
        ├── contact.page.html
        ├── details.page.html
        ├── dummy.page.html
        ├── home.page.html
        ├── test.page.html
        └── verify.page.html
├── .air.toml
├── .gitignore
├── data.json
├── go.mod
├── go.sum
├── main.go
└── README.md
```
## Code Style
Adhere to the following code style guidelines -> [ good practices](https://learn.zone01kisumu.ke/git/root/public/src/branch/master/subjects/good-practices/README.md)

### Contributing

To contribute to the project, follow these steps:

### 1. Fork the Repository 

Fork the repository on GitHub to your account.

### 2. Create a Branch

Create a new branch for  your feature or bug fix.

```sh
git checkout -b feature-name
```

### 3. Make Changes 
Make your changes and commit them with a meaningful commit messages.

```sh
git add .
git commit -m "Add feature name"
```

### 4. Push Changes

Push your changes to your forked repository.

```sh
git push origin feature-name
```

### 5. Create a Pull Request

Create a pull request from your branch to the `main` branch of the original repository. Provide a clear description of your changes.

## Code of Conduct

All contributors are expected to adhere to the project's [Code of Conduct](code-of-conduct.md). Be respectful and considerate in all interactions with the community.