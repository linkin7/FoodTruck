# The Food Truck Finder Service
A simple location based service that provides a user information about nearest available Food Trucks. This service also allows a registered user to open a Food Truck service. It's hosted in Google App engine with a MySQL instance.

[link](https://foodtruckapplication-175305.appspot.com)

## Usecase

## Design Overview

The project focuses on the back-end, more specifically canonical structure of the backend components. All the componenets are written in Go, because

- Go has in built feature for concurrency and minimal server/client based communication.
- Unit testing are easier without using any external framework.
- Objects can be accessed and returned as both type and reference variables. Also because of garbage collector mechanism, no extra headache for dangling pointers.
- Go interfaces are comparatively flexible, which helps to impose dependency injection through out the whole codebase.
- I'm more comfortable with less verbose but Go style concise coding :)

All the client components communicate with server through RPC. I chose this method over REST, SOAP, because previously I mostly worked on RPC method. Also using proper RPC based framework like stubby (gRPC), it's possible to make client implementation less dependent on server implementation. 

## High Level Description 

* **[Front-End Server](https://github.com/linkin7/FoodTruck/tree/master/src/frontendserver)**: This is the topmost component of the application. This component displays information related to user registration , finding nearest food truck etc. It communicates with application server through RPC by which it puts out the results to the browser/client tier. In simple terms, it is a layer which users can access directly (such as a web page, mobile app). During intialization this server registers all the url path with the corresponding handler function through [handler package](https://github.com/linkin7/FoodTruck/tree/master/src/frontendserver/handler).

* **[Application Server](https://github.com/linkin7/FoodTruck/tree/master/src/applicationserver)**: This is the middle layer of the architecture. This componenet interacts with component related with user data and food truck data, and exposes processed request through RPC interface. Apart from acessing data from internel components, ideally this components doesn't do any processing functionaity.
![alt text](https://github.com/linkin7/FoodTruck/blob/master/diagram.jpg)
* **[User Data Manager](https://github.com/linkin7/FoodTruck/tree/master/src/userdb)**: This components implements [UserDbManager](https://github.com/linkin7/FoodTruck/blame/master/src/common/data_manager_interface.go#L17) interface and wraps any cloud data storage for registered [user data](https://github.com/linkin7/FoodTruck/blob/master/sql.txt#L3). Because of consistency, ideally it should use any RDMS like MySQL. Apart from MySQL implementation current directory contains another in memory based mock implementation for integration testing.

* **[Food Truck Data Manager](https://github.com/linkin7/FoodTruck/tree/master/src/foodtruckdb)**: This components implements [FoodTruckDbManager](https://github.com/linkin7/FoodTruck/blob/master/src/common/data_manager_interface.go#L6) interface and wraps any cloud data storage for location based [food truck data](https://github.com/linkin7/FoodTruck/blob/master/sql.txt#L21). For time constraint currently this data storage implemented by MySQL. Apart from that, this directory also contains an in memory based mock implementation for testing.

* **[Data Container](https://github.com/linkin7/FoodTruck/tree/master/src/datacontainer/)**: Data container holds location based data in memory and provides API to insert, delete and query nearest data effciently. It should implement [DataContainer](https://github.com/linkin7/FoodTruck/blob/master/src/common/data_manager_interface.go#L25) interface and store the data in RAM. Though it should wrap a QuadTree based solution, because of simplicity and time constraint, current codebase uses a array based mock implementation.

* **Cluster Manager**: This a Map Reduce authoring batch job. By design, this job imports all the food truck data from data store, finds some uniformly distributed cluster and assigns latest cluster number to all the food trucksbaed on location. Because of the periodic mutation of food truck, it's better to store food truck data in a version controlled NoSQL database like BigTable. By modifying K means algorithm, it should be feasible to find uniformly distributed cluster centroids. Implementation should be done with Google flume like project Apache crunch. But having lack of knowledge with Apache crunch interface, currently this pipeline isn't implemented. 

* **Food Truck Data Server**: This component processes any request for food truck data. It holds a reference of user data manager, food truck data manager and data container. By design, it should store food truck data of few clusters data in data container and serves the query using data container. It's a read heavy component and processes most of the data, so it needs to be scaled widely. This component can also exploit geo location of datacentre. Because most of the request are served from nearest datacentre, each invidual server will only hold the data of nearest clusters and can store in the in-memory data container. For simplicity current implementation assigns all the food truck data to cluster 0.




