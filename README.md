
# EvoLoadBalancer

It's basically a load-balancer which has some advance feature with improved load balancing algorithm and can check the health status of the servers periodically.


## Special Features
1.Healthcheck struct indicates health of the servers here

2.The simpleServer struct now includes the healthCheck field. The newServer function is modified to accept the health check URL, interval, and timeout as parameters.

3. A new startHealthCheck method is added to the simpleServer struct, which starts a goroutine to periodically perform health checks. The checkHealth method sends an HTTP request to the health check URL and updates the server's health status based on the response.

4. the load balancer will periodically check the health of each backend server based on the provided health check URL, interval, and timeout. If a server fails the health check, it will be marked as unhealthy and will not receive new requests until it becomes healthy again.

5. The code has some improved algorithms for loadbalancing such as round-robin count, weighted round-robin count, least-connection





## Run Locally

Clone the project

```bash
  git clone https://github.com/RickDeb2004/EvoLoadBalancer
```

Go to the project directory

```bash
  cd go-loadbalancer
```

Install dependencies

```bash
  go mod init
```

Start the server

```bash
  go run main.go
```


## Tech Stack

**Go lang, Concept Of Computer Science Fundamentals**


## Advantages
1.The load balancer also keeps track of the responses from the backend servers and may cache them to improve performance. It can serve cached responses for subsequent identical requests, reducing the load on the backend servers.

2.As the load on the system changes, the load balancer dynamically adjusts the distribution of incoming requests among the backend servers to ensure optimal utilization of resources and maintain high availability.

3.This setup allows you to horizontally scale your web application by adding more backend servers as the traffic grows, and the load balancer takes care of distributing the traffic across them evenly.

Overall, the load balancer helps improve the scalability, performance, and availability of your web application by efficiently distributing the workload across multiple backend servers.




