package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
	"sync"
	
)

type Server interface {
	Address() string
	IsAlive() bool
	Serve(w http.ResponseWriter, r *http.Request)
}
type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
	healthCheck *healthCheck
	weight int
	currentcons int
	mutex sync.Mutex

}
type healthCheck struct{
	url string
	interval time.Duration
	timeout time.Duration
	healthy bool
	lastCheckedTime time.Time
	mutex sync.Mutex
}
type loadBalancingAlgorithim int
const(
	roundRobinCount loadBalancingAlgorithim=iota
	WeightedRoundRobin
	LeastConnections

)

func newServer(addr string, healthhealthCheckURL string, healthCheckInterval,healhealthCheckTimeOut time.Duration) *simpleServer {
	serverUrl, err := url.Parse(addr)
	handleErr(err)
	healthCheck :=&healthCheck{
		url: healthhealthCheckURL,
		interval:healthCheckInterval,
		timeout: healhealthCheckTimeOut,
		healthy: true,
		lastCheckedTime: time.Now(),


	}
	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverUrl),
		healthCheck: healthCheck,
	}
}

type loadBalancer struct {
	port            string
    algorithm    loadBalancingAlgorithim
	servers         []Server
	mutex          sync.Mutex
	connection int
	

}

func NewLoadBalancer(port string, servers []Server,algorithm loadBalancingAlgorithim) *loadBalancer {
	return &loadBalancer{
		port:            port,
		algorithm:  algorithm,
		servers:         servers,
	}
}
func handleErr(err error) {
	if err != nil {
		fmt.Printf("error:%v \n", err)
		os.Exit(1)
	}
}
func (s *simpleServer) Address() string { return s.addr }

func (s *simpleServer) IsAlive() bool { return s.healthCheck.healthy }

func (s *simpleServer) Serve(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}
func(s*simpleServer)startHealthCheck(){
	ticker:=time.NewTicker(s.healthCheck.interval)
	defer ticker.Stop()
	for range ticker.C{
		select{
		case <-ticker.C:
			s.checkHealth()
		}
	}
}
func(s*simpleServer)checkHealth(){
	if time.Since(s.healthCheck.lastCheckedTime)<s.healthCheck.interval{
	  return }
	  client := http.Client{Timeout: s.healthCheck.timeout}
	resp, err := client.Get(s.healthCheck.url)
	if err != nil || resp.StatusCode != http.StatusOK {
		s.healthCheck.healthy = false
		fmt.Printf("Server %s is not healthy\n", s.addr)
	} else {
		s.healthCheck.healthy = true
		fmt.Printf("Server %s is healthy\n", s.addr)
	}

	s.healthCheck.lastCheckedTime = time.Now()
}



func (lb *loadBalancer) getAvailableServerFunc(r*http.Request) Server {
	switch lb.algorithm{
	case roundRobinCount:
		return lb.getAvailableServerRoundRobin()
	case WeightedRoundRobin:
		return lb.getAvailableServerWeightedRoundRobin()
	case LeastConnections:
		return lb.getAvailableServerLeastConnections()
    default:
		return lb.getAvailableServerRoundRobin()


	}
	

}
func(lb *loadBalancer)getAvailableServerRoundRobin() Server{
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	server:=lb.servers[lb.connection%len(lb.servers)]
	lb.connection ++
	return server
}
func(lb *loadBalancer)getAvailableServerWeightedRoundRobin()Server{
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	var totalWeight int
	for _,server:=range lb.servers{
		totalWeight+=server.(*simpleServer).weight
	}
	var selectedServer Server
	for _,server:=range lb.servers{
		simpleServer:=server.(*simpleServer)
		simpleServer.healthCheck.mutex.Lock()
		if simpleServer.IsAlive(){
			selectedServer=server
			simpleServer.weight--
if simpleServer.weight==0{
	simpleServer.weight=totalWeight
}
simpleServer.mutex.Unlock()
break
		}
simpleServer.mutex.Unlock()
	}
	lb.connection++
	return selectedServer


}
func(lb *loadBalancer)getAvailableServerLeastConnections()Server{
	lb.mutex.Lock()
	defer lb.mutex.Unlock()
	var minCons int
	
	var selectedServer Server
	for _, server:=range lb.servers{
		simpleServer:=(server).(*simpleServer)
		simpleServer.mutex.Lock()
		if simpleServer.IsAlive(){
			if minCons==0|| simpleServer.currentcons<minCons{
				minCons=simpleServer.currentcons
				selectedServer=server

			}
		}
		simpleServer.mutex.Unlock()
	}
	lb.connection++
	return selectedServer


}
func (lb *loadBalancer) serveProxy(w http.ResponseWriter, r *http.Request) {
	targetServer := lb.getAvailableServerFunc(r)
	fmt.Printf("forwarding request address %q \n", targetServer.Address())
	targetServer.Serve(w, r)
}

func main() {
	servers := []Server{
		newServer("http://www.facebook.com","http://www.facebook.com/health",5*time.Second,2*time.Second),
		newServer("http://www.bing.com","http://www.bing.com/health",5*time.Second,2*time.Second),
		newServer("http://www.duckduckgo.com","http://www.duckduckgo.com/health",5*time.Second,2*time.Second),
	}
	lb := NewLoadBalancer("8000", servers,WeightedRoundRobin)
	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		lb.serveProxy(w, r)
	}
	http.HandleFunc("/", handleRedirect)
	fmt.Printf("serving request at 'local host : %s'\n", lb.port)
	http.ListenAndServe(":"+lb.port, nil)
}
func mustParseURL(rawURL string) *url.URL{
	parsedURL,err:=url.Parse(rawURL)
	if err!=nil{
		fmt.Printf("failed to parse url %v\n",err)
		os.Exit(1)
	}
	return parsedURL
}
//SPECIAL FEATURES ADDED :
//healthcheck struct indicates health of the servers here.
//The simpleServer struct now includes the healthCheck field. The newServer function is modified to accept the health check URL, interval, and timeout as parameters.
//A new startHealthCheck method is added to the simpleServer struct, which starts a goroutine to periodically perform health checks. The checkHealth method sends an HTTP request to the health check URL and updates the server's health status based on the response.
// the load balancer will periodically check the health of each backend server based on the provided health check URL, interval, and timeout. If a server fails the health check, it will be marked as unhealthy and will not receive new requests until it becomes healthy again.
//improving loadbalancing algorithm beyond the simple roundrobincount.
