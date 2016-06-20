package eago

type Distributor struct {
	status   string
	Requests RequestChan
	stop     chan struct{}
}

func NewDistributor() *Distributor {
	res := &Distributor{
		status:   STOP,
		Requests: make(RequestChan),
	}
	return res
}

func (r *Distributor) Run() {
	Log.Println("Distributor is running...")
	r.status = RUNNING
	for {
		select {
		case <-r.stop:
			Log.Println("the Distributor is stop!")
			r.stop = nil
			r.status = STOP
			return
		case reqs := <-r.Requests:
			// no goroutine, this is to reuse one connection to do the rpc, to save resources
			r.handle(reqs)
		}
	}
}

func (r *Distributor) handle(reqs []*UrlRequest) {
	for _, req := range reqs {
		if req.cookieJar == 0 {
			req.node = GetClusterInstance().GetNode(req.url)
		}
		// check if the req's node name is local node
		if GetNodeInstance().Info.NodeName == req.node {
			Log.Println("distribute the url to self: ", req.url)
			GetNodeInstance().crawl.req.push(req)
		} else {
			Log.Println("[RPC] distribute the url to ", req.node, ". url:", req.url)
			GetNodeInstance().rpc.Distribute(req)
		}
	}
}

func (r *Distributor) Stop() {
	defer func() {
		if err := recover(); err != nil {
			Error.Println(err)
		}
	}()
	close(r.stop)
}

// this func will restart the Reporter
func (r *Distributor) Restart() {
	r.stop = make(chan struct{})
	go r.Run()
}

func (r *Distributor) Status() string {
	return r.status
}