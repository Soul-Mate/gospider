package node

//import (
//	"strings"
//	"strconv"
//	"log"
//	"github.com/Soul-Mate/gospider/core"
//	"net/http"
//)
//
//type ClusterNode struct {
//	IP       string
//	Port     int
//	Name     string
//	Worker   int
//	Address  string
//	Engine   *core.Engine
//	download *core.Downloader
//}
//
//func NewClusterNode(name, address string, worker int, engine *core.Engine) *ClusterNode {
//	node := &ClusterNode{}
//	ss := strings.Split(address, ":")
//	node.Name = name
//	node.IP = ss[0]
//	if p, err := strconv.Atoi(ss[1]); err != nil {
//		log.Fatalf("The address[%s] illegality", address)
//	} else {
//		node.Port = p
//	}
//	node.Address = address
//	node.Worker = worker
//	node.Engine = engine
//	node.download = core.NewDownloader()
//	return node
//}
//
//func (n *ClusterNode) Run() {
//	// register service
//	n.download.Start()
//
//	n.callDownloader()
//
//	n.do()
//
//	log.Printf("http server start: %s", n.Address)
//
//	http.ListenAndServe(n.Address, nil)
//}
//
//func (n *ClusterNode) do() {
//	for {
//		res := n.download.GetResponse()
//		if spider, ok := n.Engine.Spiders[res.SpiderName]; ok {
//			if callback, ok := spider.Responses[res.CallBackName]; ok {
//				req := <-callback(spider, res)
//				// send master
//				println("send req to master:", req)
//			}
//		}
//	}
//}
//
//func (n *ClusterNode) callDownloader() {
//	go func() {
//		for {
//			if req := n.Engine.Schedule.PopRequest(); req != nil {
//				n.download.Commit(req)
//			}
//		}
//	}()
//}
