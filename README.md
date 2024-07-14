# middleware
Middleware logging is a technique used in software development, particularly in web and microservices applications, to log important information about incoming requests, outgoing responses, and the operations performed by the application.
- [core-go/middleware](https://github.com/core-go/middleware) is designed to integrate with middleware logging seamlessly for existing Go libraries: [Echo](https://github.com/labstack/echo), [Gin](https://github.com/gin-gonic/gin), or net/http ([Gorilla mux](https://github.com/gorilla/mux), [Go-chi](https://github.com/go-chi/chi)), with any logging libraries ([zap](https://pkg.go.dev/go.uber.org/zap), [logrus](https://github.com/sirupsen/logrus)), to log request headers, request body, response status code, body content, response time, and size
- Especially, [core-go/middleware](https://github.com/core-go/middleware) supported to encrypt sensitive data, which is useful for Financial Transactions (to comply with <b>PCI-DSS</b> standards) and Healthcare (to comply with <b>HIPAA</b> regulations)

### A typical micro service
- When you zoom one micro service, the flow is as below, and you can see "middleware" in the full picture:
  ![A typical micro service](https://cdn-images-1.medium.com/max/800/1*d9kyekAbQYBxH-C6w38XZQ.png)

## Content for logging
### Request
#### Features
- <b>Log Request Method and URL</b>: Log the HTTP method (GET, POST, etc.) and the requested URL.
- <b>Log Request Headers</b>: Option to log request headers for debugging purposes.
- <b>Log Request Body</b>: Option to log the request body (with configurable size limits to avoid logging large payloads).
#### Benefits
- <b>Debugging</b>: Helps in tracing and debugging issues by providing complete information about incoming requests.
- <b>Monitoring</b>: Provides visibility into the types of requests being received.

### Response
#### Features
- <b>Log Response Status Code</b>: Log the HTTP status code of the response.
- <b>Log Response Headers</b>: Option to log response headers.
- <b>Log Response Body</b>: Option to log the response body (with configurable size limits to avoid logging large payloads).
#### Benefits
- <b>Debugging</b>: Assists in diagnosing issues by providing complete information about the responses sent by the server.
- <b>Auditing</b>: Helps in auditing and reviewing server responses for compliance and monitoring purposes.

### Response Time
#### Features
- <b>Log Response Time</b>: Calculate and log the time taken to process each request.
#### Benefits
- <b>Performance Monitoring</b>: Helps in identifying slow requests and performance bottlenecks.
- <b>Optimization</b>: Provides data to optimize and improve server response times.

### Response Size
#### Features
- <b>Log Response Size</b>: Log the size of the response payload in bytes.
#### Benefits
- <b>Bandwidth Monitoring</b>: Helps in monitoring and managing bandwidth usage.
- <b>Optimization</b>: Provides insights into the response sizes to optimize payloads and improve performance.

## Features
### Middleware Integration
#### Features
- <b>Middleware Function</b>: Designed to integrate seamlessly with existing Go libraries: [Echo](https://github.com/labstack/echo), [Gin](https://github.com/gin-gonic/gin), or net/http ([Gorilla mux](https://github.com/gorilla/mux), [Go-chi](https://github.com/go-chi/chi)).
  - Sample for [Echo](https://github.com/labstack/echo) is at [go-sql-echo-sample](https://github.com/go-tutorials/go-sql-echo-sample)
  - Sample for [Gin](https://github.com/gin-gonic/gin) is at [go-sql-gin-sample](https://github.com/go-tutorials/go-sql-gin-sample)
  - Sample for [Gorilla mux](https://github.com/gorilla/mux) is at [go-sql-sample](https://github.com/go-tutorials/go-sql-sample) 
- <b>Context Handling</b>: Pass context to handle request-specific data throughout the middleware chain.
#### Benefits
- <b>Ease of Use</b>: Simplifies the integration of logging into existing web applications.
- <b>Consistency</b>: Ensures consistent logging across different parts of the application.

### Logging Libraries Integration
- Do not depend on any logging libraries.
- Already supported to integrate with [zap](https://pkg.go.dev/go.uber.org/zap), [logrus](https://github.com/sirupsen/logrus)
- Can be integrated with any logging library.

### Sensitive Data Encryption
#### Features
- Mask/Encrypt sensitive data in the request and response bodies.
- Sensitive Data Identification: identify and encrypt specific fields in JSON payloads.

#### Benefits:
- <b>Security</b>: Protects sensitive information from being exposed in logs.
- <b>Compliance</b>: Helps meet security and compliance requirements by safeguarding sensitive data.
- <b>Ease of Use</b>: Simplifies the integration of encryption/masking into any existing applications.
- <b>Consistency</b>: Ensures that sensitive data is consistently encrypted or masked across all logged requests and responses

#### Samples:
- Sample for [Echo](https://github.com/labstack/echo) is at [go-sql-echo-sample](https://github.com/go-tutorials/go-sql-echo-sample)
- Sample for [Gin](https://github.com/gin-gonic/gin) is at [go-sql-gin-sample](https://github.com/go-tutorials/go-sql-gin-sample)
- Sample for [Gorilla mux](https://github.com/gorilla/mux) is at [go-sql-sample](https://github.com/go-tutorials/go-sql-sample)

### Enable/Disable Logging
#### Features
- <b>Enable/Disable Logging</b>: Allow users to turn on or off logging for requests, responses, headers, and bodies independently.
- <b>Logging Levels</b>: Support different logging levels (e.g., INFO, DEBUG, ERROR) to control the verbosity of logs.
#### Benefits
- <b>Flexibility</b>: Provides users with the flexibility to configure logging based on their needs and environment.
- <b>Efficiency</b>: Reduces overhead by allowing selective logging, especially in production environments.

### Asynchronous Logging
#### Features
- <b>Non-Blocking Logs</b>: Implement asynchronous logging to ensure that logging does not block request processing.
- <b>Log Buffering</b>: Use buffering to improve logging performance and reduce latency.
#### Benefits:
- <b>Performance</b>: Improves the overall performance of the application by reducing logging overhead.
- <b>Scalability</b>: Allows the application to handle high-throughput logging without performance degradation.


## Use Cases of Sensitive Data Encryption
### Financial Transactions
- <b>Benefit</b>: Encrypting sensitive financial data, such as credit card numbers and transaction details, helps comply with PCI-DSS standards and secures financial transactions from exposure in logs.
### Healthcare
- <b>Benefit</b>: Encrypting patient data such as medical records and health information in logs ensures compliance with HIPAA regulations and protects patient privacy.
### E-commerce
- <b>Benefit</b>: Protecting customer information, such as addresses and payment details, enhances customer trust and protects the e-commerce platform from potential data breaches.

## Benefits of Middleware Logging
#### Debugging and Troubleshooting
- Provides detailed logs that help developers debug and troubleshoot issues in the application by tracing the flow of requests and responses.
#### Monitoring and Alerting
- Enables monitoring of application performance and behavior, allowing for real-time alerting on errors, slow responses, and unusual activity.
#### Performance Optimization
- Logs performance metrics that can be analyzed to identify bottlenecks, optimize resource usage, and improve overall application performance.
#### Security and Compliance
- Helps in tracking access and usage patterns, detecting security incidents, and complying with regulatory requirements by logging relevant information.
#### Auditing
- Provides an audit trail of user actions and system operations, which is essential for security audits and forensic analysis.

## Conclusion
Middleware logging is a critical aspect of building robust, maintainable, and secure applications, providing valuable insights and aiding in the continuous improvement of the software.

## Installation
Please make sure to initialize a Go module before installing core-go/middleware:

```shell
go get -u github.com/core-go/middleware
```

Import:
```go
import "github.com/core-go/middleware"
```

## Appendix
### Microservice Architect
![Microservice Architect](https://cdn-images-1.medium.com/max/800/1*vKeePO_UC73i7tfymSmYNA.png)

### A typical micro service
- When you zoom one micro service, the flow is as below, and you can see "middleware" in the full picture:
  ![A typical micro service](https://cdn-images-1.medium.com/max/800/1*d9kyekAbQYBxH-C6w38XZQ.png)

### Cross-cutting concerns
- "middleware" in the full picture of cross-cutting concerns
  ![cross-cutting concerns](https://cdn-images-1.medium.com/max/800/1*y088T4NoJNrL9sqrKeSyqw.png)
