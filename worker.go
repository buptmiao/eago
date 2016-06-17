package crawler



const(
	STOP     =   "stop"
	RUNNING  =   "running"
)

// the Worker defines these three method:Run, Stop, Restart.
// Actually, Fetcher, Extractor, Reporter and Distributor
// implement this interface. Every worker has a dead loop to
// do a specific task in one go routine unless Stop is invoked
type Worker interface {
	Run()
	Stop()
	Restart()
}
