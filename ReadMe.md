
### Project Overview
The **Go Web Analyzer** is a tool designed to analyze web pages for specific elements such as links, headings, and other metadata. It provides insights into the structure and content of web pages, making it useful for developers, SEO analysts, and content creators. The application emphasizes robust error handling, logging, and performance monitoring.


### Usage

- Once you start up the application using `make clean run` or containerizing the app using the `make docker-run`
  The application can be run in the browser with `http://localhost:8080` and then UI will make things easier to explore the app.

- Other than using the browser. You can use a `POST http://localhost:8080/analyze` with following request body  
  gives you the analyze response in a JSON format.

  ```
    {
        "url":"http://example.com"
    }
  ```

  URL is valid only when url starts as `http://` or `https://`

- Inorder to watch the metrics `GET http://localhost:8080/metrics` endpoint can be used.

### Prerequisites
- **Go Programming Language**: Ensure Go is installed on your system. 
- **Docker with Docker Desktop** (optional): For containerized deployment. 

### Technologies Used
#### Backend (BE)
- **Language**: Go
- **Frameworks**: Standard Go libraries for HTTP and HTML content extraction

#### Frontend (FE)
- Frontend is a simple Index HTML created for the purpose of Presentation

#### DevOps
- **Monitoring**: Prometheus
- **Containerization**: Docker


### External Dependencies
- Using [golang.org/x/net/html](https://pkg.go.dev/golang.org/x/net/html) library for HTML Parsing. 

### Setup Instructions
To install and run the project, follow these steps:

1. **Clone the Repository**:
    ```bash
    git clone https://github.com/Jawadh-Salih/go-web-analyzer.git
    cd go-web-analyzer
    ```

2. **Install Dependencies**:
    Ensure all required Go modules are installed:
    ```bash
    go mod tidy
    ```

3. **Build the Application**:
    Use the provided `Makefile` to build the project:
    ```bash
    make clean build
    ```

4. **Run the Application**:
    Start the application using:
    ```bash
    make clean run
    ```

5. **Run Tests**:
    Execute the test suite to ensure everything is working correctly:
    ```bash
    make test
    ```

6. **Test Coverage**:
    Test covergae can be checked in following ways:
    - To see the coverage in the terminal
    ```bash
    make coverage
    ```
    - To see the coverage in the HTML
    ```bash
    make coverage-html-open
    ```


For additional details, refer to the `Makefile` and project documentation.
present


### Challenges faced and the approaches took to overcome

Parsing HTML content is not straight forward and I had rely on a library (mentioned above) and use a recursive implementation to Traverse through each HTML nodes for different functionalities.

I was having a performance issue when ran the analyzer logic sequencially for each different requirement. So having concurrent implementation optimized the logic immensely. For that I have used waitgroups and channels with a Defined struct to acheive a smooth parallelism. With this I managed to reduce the analyzer time from 9s to 6s.

Further to this, I had another performance issue on Extracting the Links of the url where the links required to be tested for their accesbility. For that we need to call each of the links and see whether they give a 2XX response. Since the urls can be independently called I used worker pool and buffered channel to acheive concurrency and hence I was able to reduce the logic from ~6.5s to ~1.5s (this was tested for the url https://github.com/login ).

I was using http.Get to check the URL accessibility. But it was not performant enough. So I used http.Head to make it faster. 


### Future Improvements
- Enhance error handling with Error codes to identify errors better.
- Enhance the URL accessibility check by adding support for retry logic and handling rate-limiting scenarios.
- Introduce a configuration file to allow users to customize the analyzer's behavior, such as setting timeouts or enabling/disabling specific features.
- Improve the user interface by creating a more interactive and visually appealing frontend.
- Integrate with third-party APIs to provide additional insights, such as SEO scores or content readability analysis.
- Provide multilingual support for analyzing web pages in different languages.
- Develop a plugin system to allow users to extend the application's functionality with custom modules.
- Explore the use of machine learning to identify patterns or anomalies in web page structures.
- Create a CLI version of the application for easier integration into automated workflows.
- Add support for exporting analysis results in various formats, such as JSON, CSV, or PDF.
- Introduce CI/CD pipeline for the project

### Additional Notes
- Ensure all dependencies are installed before running the project.
- Refer to the respective repository URLs for detailed setup instructions.
- The application is designed with modularity in mind, allowing easy integration of additional features in the future.
- Logging and error handling are critical components, ensuring smooth debugging and monitoring.
 