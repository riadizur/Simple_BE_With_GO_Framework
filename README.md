Tittle      : Stechoq Backend Framework with GO

Author      : @riadizur

Description : 
This simple backend framework built specifically for Stechoq backend developers, designed to streamline the process of creating web applications. With This simple backend , Stechoq developers can effortlessly set up HTTP REST APIs, WebSocket connections, SQLite3 database integration, process logging, and environment variable configuration.

Key Features:
1. HTTP REST API: Stechoq simplifies the creation of RESTful APIs, allowing developers to define endpoints, handle requests, and manage responses with ease. Whether you're building a simple CRUD application or a complex microservices architecture, Stechoq provides the tools you need to get the job done efficiently.
2. WebSocket Support: Real-time communication is essential for many modern web applications, and Stechoq makes it simple to implement WebSocket functionality. Whether you're building a chat application, a multiplayer game, or a live data visualization tool, Stechoq's WebSocket support has you covered.
3. SQLite3 Database Integration: Database access is a fundamental aspect of many backend applications, and Stechoq seamlessly integrates with SQLite3 to provide efficient data storage and retrieval. Whether you're working with structured data or unstructured content, Stechoq's database integration ensures optimal performance and reliability.
4. Process Logging: Monitoring and debugging are crucial aspects of backend development, and Stechoq includes robust logging functionality to help you track the execution of your application. From error messages to performance metrics, Stechoq's logging capabilities ensure that you have the insights you need to keep your application running smoothly.
5. Environment Variable Configuration: Configuration management is essential for deploying applications across different environments, and Stechoq simplifies this process with support for .env files. By centralizing your configuration settings in a single file, Stechoq makes it easy to manage your application's behavior across development, staging, and production environments.
Whether you're a seasoned Go developer or just getting started with the language, this framework as toolkit for building robust backend applications. With its intuitive API and documentation to empowers developers to focus on building great software without getting bogged down by infrastructure concerns.

The Step of Instruction :
1. Prepare workspace with running shell install.sh to providing Folder and Initiation GO Workspaces.
CODE STRUCTURE Will be provided :

CODE STRUCTURE
 ![alt text](https://github.com/riadizur/Simple_BE_With_GO_Framework/blob/main/code_structure.png)

kubota-gasoline-api/
├── cmd/
│   └── main.go
├── db/
│   └── kubota_gasoline.db
├── internal/
│   ├── database/
│   │   └── database.go
│   ├── handlers/
│   │   ├── shift.go
│   │   ├── websocket.go
│   └── models/
│       └── shift.go
├── log/
│   └── process.log
├── go.mod
└── go.sum
└── (application_build)

2. Let's cook your code on CMD and Internal Folder.

3. Run build.sh to build application on your pc inveronment or Run build-win.sh to cross-compile the application that run on Windows.

Thank You.
