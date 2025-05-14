
- Ensure crucial Tests are covered. 
TODO
    - Extracting Links test
    - Extracting Headings test
    - Ensure logging in place espcially with errors. 
    - can we pass the logger into analyzer?
    - Include a prometheous
    - Check Error Handling

Fill the readme


### Project Overview
The **Go Web Analyzer** is a tool designed to analyze web pages for specific elements such as links, headings, and other metadata. It provides insights into the structure and content of web pages, making it useful for developers, SEO analysts, and content creators. The application emphasizes robust error handling, logging, and performance monitoring.

### Prerequisites
- **Go Programming Language**: Ensure Go is installed on your system. [Download Go](https://go.dev/dl/)
- **Prometheus**: For monitoring and metrics collection. [Prometheus Installation Guide](https://prometheus.io/docs/introduction/install/)
- **Docker** (optional): For containerized deployment. [Docker Installation Guide](https://docs.docker.com/get-docker/)

### Technologies Used
#### Backend (BE)
- **Language**: Go
- **Frameworks**: Standard Go libraries for HTTP and HTML content extraction

#### Frontend (FE)
- Frontend is a simple Index HTML created for the purpose of Presentation

#### DevOps
- **CI/CD**: GitHub Actions for automated testing and deployment
- **Monitoring**: Prometheus and Grafana
- **Containerization**: Docker


### Additional Notes
- Ensure all dependencies are installed before running the project.
- Refer to the respective repository URLs for detailed setup instructions.
- The application is designed with modularity in mind, allowing easy integration of additional features in the future.
- Logging and error handling are critical components, ensuring smooth debugging and monitoring.


### External Dependencies
- Using golang.org/x/net/html library for HTML Parsing. 

### Setup Instructions

To install and run the project, follow these steps:

1. **Clone the Repository**:
    ```bash
    git clone https://github.com/your-repo/go-web-analyzer-lucytech.git
    cd go-web-analyzer-lucytech
    ```

2. **Install Dependencies**:
    Ensure all required Go modules are installed:
    ```bash
    go mod tidy
    ```

3. **Build the Application**:
    Use the provided `Makefile` to build the project:
    ```bash
    make build
    ```

4. **Run the Application**:
    Start the application using:
    ```bash
    make run
    ```

5. **Run Tests**:
    Execute the test suite to ensure everything is working correctly:
    ```bash
    make test
    ```

8. **Access the Application**:
    Open your browser and navigate to `http://localhost:8080` (or the configured port) to access the application.

For additional details, refer to the `Makefile` and project documentation.
present


● Mention the usage of the App with main functionalities and their role in the
Application (Eg: authentication, logging, error handling)

### Challenges you faced and the approaches you took to overcome

Parsing HTML content is not straight forward and I had rely on a library (mentioned above) and use a recursive implementation to Traverse through each HTML nodes for different functionalities.

I was having a performance issue when ran the analyzer logic sequencially for each different requirement. So having concurrent implementation optimized the logic immensely. For that I have used waitgroups and go channels with a Defined struct to acheive a smooth parallelism. With this I managed to reduce the analyzer time from 9s to 6s.

Further to this, I had another performance issue on Extracting the Links of the url where the links required to be tested for their accesbility. For that we need to call each of the links and see whether they give a 2XX response. Since the urls can be independently called I used worker pool and buffered channel to acheive parallelism and hence I was able to reduce the logic from 6s to 1s. 


### Future Improvements

- Implement a caching mechanism to store previously analyzed results, reducing redundant processing for frequently analyzed URLs.
- Enhance the URL accessibility check by adding support for retry logic and handling rate-limiting scenarios.
- Introduce a configuration file to allow users to customize the analyzer's behavior, such as setting timeouts or enabling/disabling specific features.
- Improve the user interface by creating a more interactive and visually appealing frontend.
- Integrate with third-party APIs to provide additional insights, such as SEO scores or content readability analysis.
- Provide multilingual support for analyzing web pages in different languages.
- Develop a plugin system to allow users to extend the application's functionality with custom modules.
- Explore the use of machine learning to identify patterns or anomalies in web page structures.
- Create a CLI version of the application for easier integration into automated workflows.
- Add support for exporting analysis results in various formats, such as JSON, CSV, or PDF.

