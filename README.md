
# eago

## Installation
    go get github.com/gin-gonic/gin
    go get gopkg.in/redis.v3
    go get github.com/BurntSushi/toml

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


