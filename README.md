# eago
An easy distribute and restful crawler framework, this can be

## Installation
Before install eago, you should install these dependence:

    go get github.com/gin-gonic/gin
    go get gopkg.in/redis.v3
    go get github.com/BurntSushi/toml

Install:

    go get github.com/buptmiao/eago
    
## Usage
```go
import (
	"github.com/buptmiao/eago"
)

func main() {

	eago.LoadConfig()
	node := eago.GetNodeInstance()
	cluster := eago.GetClusterInstance()

	eago.NewRpcServer().Start()
	// Descover will Block the execution, until a master node
	// is found, or become master itself.
	cluster.Discover()
	// start the Http Server
	eago.NewHttpServer(node).Serve()
}
```

#### Filter
Use Redis to filter the duplicate urls


