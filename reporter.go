package eago

type Reporter struct {
	status string
	stop chan struct{}
	pop RequestChan
}

func NewReporter(pop RequestChan) *Reporter{
	res := &Reporter{
		status : STOP,
		pop : pop,
	}
	return res
}

func (r *Reporter)Run() {
	Log.Println("Reporter is running...")
	r.status = RUNNING
	for{
		select {
		case <-r.stop:
			Log.Println("the Reporter is stop!")
			r.stop = nil
			r.status = STOP
			return
		case reqs := <-r.pop:
			// no goroutine, this is to reuse one connection to do the rpc, to save resources
			r.handle(reqs)
		}
	}
}

func (r *Reporter)handle(reqs []*UrlRequest) {
	for _, req := range reqs {
		if GetNodeInstance().IsMaster() {
			Log.Println("report the url to self: ", req.url)
			GetClusterInstance().PushRequest(req)
		} else {
			Log.Println("[RPC] report the url to master: ", req.url)
			if err := GetNodeInstance().rpc.ReportRequest(req); err != nil {
				Error.Println(err)
			}
		}
	}
}

func (r *Reporter)Stop() {
	defer func(){
		if err := recover(); err != nil {
			Error.Println(err)
		}
	}()
	close(r.stop)
}

// this func will restart the Reporter
func (r *Reporter) Restart() {
	r.stop = make(chan struct{})
	go r.Run()
}


