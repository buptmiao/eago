package crawler

import (
	"time"
	"sync/atomic"
	"fmt"
)
// CrawlerStatistic demonstrate the crawler's basic info
// used by restful api to monitor current state of the crawler
type CrawlerStatistic struct{
	Name  				string			`json:"CrawlerName"`
	Running             bool    		`json:"Running"`
	CrawledUrlsCount  	uint64			`json:"CrawledUrlsCount"`
	TotalCount			uint64			`json:"TotalUrlsCount"`
	ToCrawledUrlsCount	uint64          `json:"ToCrawledUrlsCount"`
	BeginAt 			string   		`json:"Begin At"`
	Elapse     			string   		`json:"Elapse"`
}

// The statistic information of the cluster, it will be updated by
// the master node, when the slavers need it, they must call the
// Rpc SyncStatistic to sync this.
type Statistic struct{
	ClusterName  string              `json:"ClusterName"`
	NodeNum int                      `json:"NodeNumber"`
	Master  *NodeInfo                `json:"Master"`
	Slavers []*NodeInfo              `json:"slavers"`
	Crawler CrawlerStatistic         `json:"CrawlerStatistics"`
	Message string                   `json:"Message"`
}

const(
	Message = "You know, for data"
)
var Stat = NewStatistic()

func NewStatistic() *Statistic {
	res := &Statistic{
		ClusterName:Configs.ClusterName,
		NodeNum:1,
		Message: Message,
		Crawler:CrawlerStatistic{
			Name:Configs.CrawlerName,
		},
	}
	return res
}
// record current time at which the crawler begin
func (s *Statistic)BeginNow() *Statistic{
	s.Crawler.Running = true
	s.Crawler.BeginAt = time.Now().Format("2006-01-02 15:04:05")
	return s
}

func (s *Statistic)Stop() *Statistic{
	s.Crawler.Running = false
	begin, err := time.ParseInLocation("2006-01-02 15:04:05", s.Crawler.BeginAt, time.Local)
	if err != nil {
		Error.Println("parse time failed, Crawler has not been started ever")
		Stat.Crawler.Elapse = ""
		return  s
	}
	elapse := float64(time.Since(begin)/1e6)/float64(1e3)
	s.Crawler.Elapse = fmt.Sprintf("%.2f secs", elapse)
	return s
}

// not used
func (s *Statistic)SetCrawlerName(name string) *Statistic{
	s.Crawler.Name = name
	return s
}

func (s *Statistic)SetMaster(Node *NodeInfo) *Statistic{
	s.Master = Node
	return s
}

func (s *Statistic)AddNode(Node *NodeInfo) *Statistic{
	s.NodeNum++
	s.Slavers = append(s.Slavers, Node)
	return s
}

func (s *Statistic)SetClusterName(name string) *Statistic{
	s.ClusterName = name
	return s
}

func (s *Statistic)AddTotalCount() {
	atomic.AddUint64(&Stat.Crawler.TotalCount, 1)
}

func (s *Statistic)AddCrawledCount() {
	atomic.AddUint64(&Stat.Crawler.CrawledUrlsCount, 1)
}

// Get the current info of the crawler cluster, this will always invoked
// by the master node
func (s *Statistic)GetStatistic() *Statistic{
	// copy one to avoid data race
	stat := *s
	// If the Crawler is not running, return stat directly.
	if !s.Crawler.Running {
		return &stat
	}

	begin, err := time.ParseInLocation("2006-01-02 15:04:05", stat.Crawler.BeginAt, time.Local)
	if err != nil {
		Error.Println("parse time failed, Crawler has not been started ever")
		Stat.Crawler.Elapse = ""
		return &stat
	}
	elapse := float64(time.Since(begin)/1e6)/float64(1e3)
	stat.Crawler.Elapse = fmt.Sprintf("%.2f secs", elapse)
	return &stat
}
