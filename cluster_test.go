package eago

import (
	"testing"
	"time"
)

func TestNewCluster(t *testing.T) {
	LoadTestConfig()
	NewCluster()
	AssertNotNil(DefaultCluster)
	AssertNotNil(DefaultCluster.dis)
	AssertNotNil(DefaultCluster.hash)
	AssertEqual(DefaultCluster.ClusterName == "eagles")
	AssertNotNil(DefaultCluster.Local == DefaultNode)
}

func TestGetClusterInstance(t *testing.T) {
	LoadTestConfig()
	t1 := GetClusterInstance()
	t2 := GetClusterInstance()
	AssertEqual(t1 == t2)
}

func TestCluster_AddNode(t *testing.T) {
	LoadTestConfig()
	GetClusterInstance().AddNode(&NodeInfo{
		NodeName: "testone",
		IP:       "",
		Port:     0,
	})
	for node, st := range GetClusterInstance().Nodes {
		AssertEqual(node.NodeName == "testone")
		AssertEqual(st)
	}
	v, err := GetClusterInstance().hash.Get("*****")
	AssertErrNil(err)
	AssertEqual(v == "testone")
}

//func TestCluster_PushRequest(t *testing.T) {
//	LoadTestConfig()
//	GetClusterInstance().PushRequest(&UrlRequest{
//		url: "www.github.com",
//	})
//	reqs := <-GetClusterInstance().dis.Requests
//	AssertEqual(*reqs[0] == UrlRequest{
//		url: "www.github.com",
//	})
//}

func TestCluster_BecomeMaster(t *testing.T) {
	LoadTestConfig()
	GetClusterInstance().BecomeMaster()
	AssertEqual(GetClusterInstance().Master == GetClusterInstance().Local.Info)
	GetClusterInstance().StopKeeper()
	GetClusterInstance().StopDistributor()
}

func TestCluster_Discover(t *testing.T) {
	LoadTestConfig()
	GetClusterInstance().Discover()
	AssertEqual(GetClusterInstance().Master == GetClusterInstance().Local.Info)
}

func TestCluster_GetNode(t *testing.T) {
	LoadTestConfig()
	GetClusterInstance().AddNode(&NodeInfo{
		NodeName: "testone",
		IP:       "",
		Port:     0,
	})
	AssertEqual(GetClusterInstance().GetNode("www.github.com") == "testone")
}

func TestCluster_IsMember(t *testing.T) {
	LoadTestConfig()
	GetClusterInstance().AddNode(&NodeInfo{
		NodeName: "testone",
		IP:       "",
		Port:     0,
	})
	AssertEqual(GetClusterInstance().IsMember(&NodeInfo{
		NodeName: "testone",
		IP:       "",
		Port:     0,
	}))
}

func TestCluster_StartDistributor(t *testing.T) {
	LoadTestConfig()
	GetClusterInstance().StartDistributor()
	time.Sleep(time.Second)
	AssertEqual(GetClusterInstance().dis.status == RUNNING)
}
