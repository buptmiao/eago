package eago

import (
	"fmt"
	"sync/atomic"
	"time"
)

// CrawlerStatistic demonstrate the crawler's basic info
// used by restful api to monitor current state of the crawler
type CrawlerStatistic struct {
	Name               string `json:"CrawlerName"`
	CrawledUrlsCount   uint64 `json:"CrawledUrlsCount"`
	TotalCount         uint64 `json:"TotalUrlsCount"`
	ToCrawledUrlsCount uint64 `json:"ToCrawledUrlsCount"`
}

func NewCrawlerStatistic(name string) *CrawlerStatistic {
	res := &CrawlerStatistic{
		Name:               name,
		CrawledUrlsCount:   0,
		TotalCount:         0,
		ToCrawledUrlsCount: 0,
	}
	return res
}

// The statistic information of the cluster, it will be updated by
// the master node, when the slavers need it, they must call the
// Rpc SyncStatistic to sync this.
type Statistic struct {
	ClusterName string                       `json:"ClusterName"`
	Running     bool                         `json:"Running"`
	BeginAt     string                       `json:"Begin At"`
	Elapse      string                       `json:"Elapse"`
	NodeNum     int                          `json:"NodeNumber"`
	Master      *NodeInfo                    `json:"Master"`
	Slavers     []*SlaverStatus              `json:"slavers"`
	Crawler     map[string]*CrawlerStatistic `json:"CrawlerStatistics"`
	Message     string                       `json:"Message"`
}

type SlaverStatus struct {
	*NodeInfo
	Alive bool `json:"alive"`
}

const (
	Message = "You know, for data"
)

var Stat = NewStatistic()

func NewStatistic() *Statistic {
	res := &Statistic{
		ClusterName: Configs.ClusterName,
		NodeNum:     1,
		Message:     Message,
		Crawler:     make(map[string]*CrawlerStatistic),
	}
	return res
}

// record current time at which the crawler begin
func (s *Statistic) BeginNow() *Statistic {
	s.Running = true
	s.BeginAt = time.Now().Format("2006-01-02 15:04:05")
	return s
}

func (s *Statistic) Stop() *Statistic {
	s.Running = false
	begin, err := time.ParseInLocation("2006-01-02 15:04:05", s.BeginAt, time.Local)
	if err != nil {
		Error.Println("parse time failed, Crawler has not been started ever")
		Stat.Elapse = ""
		return s
	}
	elapse := float64(time.Since(begin)/1e6) / float64(1e3)
	s.Elapse = fmt.Sprintf("%.2f secs", elapse)
	return s
}

func (s *Statistic) AddCrawlerStatistic(name string) *Statistic {
	cs := NewCrawlerStatistic(name)
	s.Crawler[name] = cs
	return s
}

func (s *Statistic) GetCrawlerStatistic(name string) *CrawlerStatistic {
	res, ok := s.Crawler[name]
	if !ok {
		panic("crawler statistic not found" + name)
		return nil
	}
	return res
}

func (s *Statistic) SetMaster(Node *NodeInfo) *Statistic {
	s.Master = Node
	return s
}

func (s *Statistic) AddNode(Node *NodeInfo) *Statistic {
	for _, node := range s.Slavers {
		if *Node == *node.NodeInfo {
			node.Alive = true
			return s
		}
	}
	s.NodeNum++
	slaverStatus := &SlaverStatus{
		NodeInfo: Node,
		Alive:    true,
	}
	s.Slavers = append(s.Slavers, slaverStatus)
	return s
}

func (s *Statistic) UpdateNodeAlive(Node *NodeInfo, v bool) *Statistic {
	for _, node := range s.Slavers {
		if *Node == *node.NodeInfo {
			node.Alive = v
			return s
		}
	}
	return s
}

func (s *Statistic) SetClusterName(name string) *Statistic {
	s.ClusterName = name
	return s
}

func (s *Statistic) AddTotalCount(name string) {
	cs := s.GetCrawlerStatistic(name)
	atomic.AddUint64(&cs.TotalCount, 1)
}

func (s *Statistic) AddCrawledCount(name string) {
	cs := s.GetCrawlerStatistic(name)
	atomic.AddUint64(&cs.CrawledUrlsCount, 1)
}

// Get the current info of the crawler cluster, this will always invoked
// by the master node
func (s *Statistic) GetStatistic() *Statistic {
	// copy one to avoid data race
	stat := *s
	// If the Crawler is not running, return stat directly.
	if !s.Running {
		return &stat
	}

	begin, err := time.ParseInLocation("2006-01-02 15:04:05", stat.BeginAt, time.Local)
	if err != nil {
		Error.Println("parse time failed, Crawler has not been started ever")
		Stat.Elapse = ""
		return &stat
	}
	elapse := float64(time.Since(begin)/1e6) / float64(1e3)
	stat.Elapse = fmt.Sprintf("%.2f secs", elapse)
	return &stat
}
