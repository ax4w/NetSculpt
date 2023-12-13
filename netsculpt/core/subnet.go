package core

import (
	"fmt"
	"math"
)

type IP struct {
	v      map[int]int
	Prefix int
}

type Data struct {
	Ips              int
	NetworkAddress   IP
	BroadcastAddress IP
	Message          string
}

var StartingIP IP

func (ip *IP) IsAllZeros() bool {
	return ip.v[0] == ip.v[1] && ip.v[1] == ip.v[2] && ip.v[2] == ip.v[3] && ip.v[3] == 0
}

func (ip *IP) ToString() string {
	str := ""
	for i := 0; i < 4; i++ {
		if i < 3 {
			str += fmt.Sprintf("%d.", ip.v[i])
		} else {
			str += fmt.Sprintf("%d", ip.v[i])
		}
	}
	str += fmt.Sprintf("/%d", ip.Prefix)
	return str
}

func (ip *IP) Show() {
	for i := 0; i < 4; i++ {
		if i < 3 {
			print(ip.v[i], ".")
		} else {
			print(ip.v[i])
		}
	}
	println("/", ip.Prefix)
}

func SetStartingIP(ipIn []int, p int) {
	m := make(map[int]int)
	for i := 0; i < 4; i++ {
		m[i] = ipIn[i]
	}
	StartingIP.v = m
	StartingIP.Prefix = p

	r := networkAddress(p)
	for i := 0; i < 4; i++ {
		StartingIP.v[i] = r[i]
	}
	//StartingIP.Show()
}

func CalculateSubnet(requiredHosts int) Data {
	var result Data
	if int(math.Pow(2, float64(32-StartingIP.Prefix))) < requiredHosts {
		return Data{Message: "Could not fit in"}
	}
	logOfHosts := logOfHosts(requiredHosts)
	result.Ips = sizeOfNetwork(logOfHosts)
	prefix := prefix(logOfHosts)
	result.NetworkAddress = toIp(networkAddress(prefix), prefix)
	result.BroadcastAddress = toIp(broadcastAddress(prefix), prefix)
	SetStartingIP(addOne(broadcastAddress(prefix)), prefix)
	return result
}

func logOfHosts(hosts int) int {
	return int(math.Ceil(math.Log2(float64(hosts + 2))))
}

func sizeOfNetwork(logOfHosts int) int {
	return int(math.Pow(2, float64(logOfHosts)))
}

func addOne(ipIn []int) []int {
	for i := 3; i >= 0; i-- {
		if ipIn[i]+1 >= 256 {
			ipIn[i] = 0
		} else {
			ipIn[i]++
			break
		}
	}
	return ipIn
}

func toIp(ipIn []int, prefix int) IP {
	m := make(map[int]int)
	for i := 0; i < 4; i++ {
		m[i] = ipIn[i]
	}
	return IP{
		v:      m,
		Prefix: prefix,
	}
}

func broadcastAddress(prefix int) []int {
	nMask := subnetMask(prefix, true)
	var res []int
	for i := 0; i < 4; i++ {
		res = append(res, StartingIP.v[i]|nMask[i])
	}

	return res
}

func networkAddress(prefix int) []int {
	nMask := subnetMask(prefix, false)
	var res []int
	for i := 0; i < 4; i++ {
		res = append(res, StartingIP.v[i]&nMask[i])
	}
	return res
}

func subnetMask(prefix int, invert bool) []int {
	curr := 0
	m := make(map[int][]int)
	for i := 0; i < 32; i++ {
		if i%8 == 0 && i != 0 {
			curr++
		}
		if i < prefix {
			if invert {
				m[curr] = append(m[curr], 0)
			} else {
				m[curr] = append(m[curr], 1)
			}
		} else {
			if invert {
				m[curr] = append(m[curr], 1)
			} else {
				m[curr] = append(m[curr], 0)
			}
		}
	}
	var r []int
	for i := 0; i < 4; i++ {
		if isAllOne(m[i]) {
			r = append(r, 255)
		} else {
			t := 0
			for i, v := range m[i] {
				if v == 1 {
					t += int(math.Pow(2, float64(7-i)))
				}
			}
			r = append(r, t)
		}
	}
	return r
}

func prefix(logOfHosts int) int {
	return 32 - logOfHosts
}
