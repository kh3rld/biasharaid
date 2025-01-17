# Usage Guide

## Basic Usage

### Running the Web Application
1. **Navigate to the Project Directory**

Open a terminal and navigate to the directory containing the Go application.
```sh
   cd  biasharaid
   ```
2. Initialized go.mod on the root of the project using the below command
``` go 
go mod init github.com/kh3rld/biasharaid
 ```
3. **Navigate to `web` directory**
    ```sh
    cd web
    ```

4. **Run the Application**
```sh
go run main.go
```

The application will start and listen on port `8080`.

## Accessing the Server

Once the server is running, you can access it via a web browser through this url:

```
http://localhost:8080/status
```
