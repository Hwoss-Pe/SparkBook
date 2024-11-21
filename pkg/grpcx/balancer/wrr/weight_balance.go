package wrr

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"sync"
)

//平滑的带权重的负载均衡

const WeightRoundRobin = "custom_weighted_round_robin"

func newBuilder() balancer.Builder {
	return base.NewBalancerBuilder(WeightRoundRobin,
		&WeightedPickerBuilder{}, base.Config{
			HealthCheck: true,
		})
}

type WeightedPickerBuilder struct {
}
type WeightedPicker struct {
	mutex sync.Mutex
	conns []*weightConn
}

func init() {
	balancer.Register(newBuilder())
}

// Pick 这里是具体的算法实现
func (w *WeightedPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	//当计算的时候没有任何节点
	if len(w.conns) == 0 {
		return balancer.PickResult{}, balancer.ErrNoSubConnAvailable
	}
	var totalWeight int
	var res *weightConn
	w.mutex.Lock()
	//算法具体逻辑，计算总权重,当前的会进行累加weight， 当是最大的时候就减掉总权重
	for _, conn := range w.conns {
		totalWeight += conn.weight
		conn.currenWeight += conn.weight
		//选当前节点最大
		if res == nil || res.currenWeight < conn.currenWeight {
			res = conn
		}
	}
	res.currenWeight -= totalWeight
	w.mutex.Unlock()
	return balancer.PickResult{
		SubConn: res.SubConn,
		Done: func(info balancer.DoneInfo) {
			//在这里可以进行调整failover相关的问题，权重调低或者直接去掉这个节点
		},
	}, nil
}

func (w *WeightedPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker {
	//ReadySCs map[balancer.SubConn]SubConnInfo
	conns := make([]*weightConn, 0, len(info.ReadySCs))
	for con, conInfo := range info.ReadySCs {
		//新版本都是写在Addr里面，老版本写在metadata里面,默认是写好的用resolver可以去实现
		weightVAl, _ := conInfo.Address.Metadata.(map[string]any)["weight"]
		// 经过注册中心的转发之后，变成了 float64，要小心这个问题
		weight := weightVAl.(float64)
		conns = append(conns, &weightConn{
			SubConn:      con,
			currenWeight: int(weight),
			weight:       int(weight),
		})
	}
	return &WeightedPicker{
		conns: conns,
	}
}

type weightConn struct {
	//初始权重
	weight int
	//当前权重
	currenWeight int
	//这个代表某个节点
	balancer.SubConn
}
