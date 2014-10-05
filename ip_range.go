package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
)

func inet_ntoa(ipnr int64) net.IP {
	var bytes [4]byte
	bytes[0] = byte(ipnr & 0xFF)
	bytes[1] = byte((ipnr >> 8) & 0xFF)
	bytes[2] = byte((ipnr >> 16) & 0xFF)
	bytes[3] = byte((ipnr >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0])
}

func inet_aton(ipnr net.IP) int64 {
	bits := strings.Split(ipnr.String(), ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64
	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)
	return sum
}

type IPRange struct {
	StartIP int64
	EndIP   int64
}

func parseIPRange(start, end string) (*IPRange, error) {
	start = strings.TrimSpace(start)
	end = strings.TrimSpace(end)

	if !strings.Contains(end, ".") {
		ss := strings.Split(start, ".")
		st := strings.Join(ss[0:3], ".")
		end = fmt.Sprintf("%s.%s", st, end)
		//		fmt.Printf("###%v  ", st)
		//		return nil, fmt.Errorf("Invalid IPRange %s-%s", start, end)
	}
	//fmt.Printf("##%s %s\n",start, end)
	si := net.ParseIP(start)
	ei := net.ParseIP(end)

	iprange := new(IPRange)
	iprange.StartIP = inet_aton(si)
	iprange.EndIP = inet_aton(ei)
	if iprange.StartIP > iprange.EndIP {
		return nil, fmt.Errorf("Invalid IPRange %s-%s", start, end)
	}
	return iprange, nil
}

func parseIPRangeFile(file string) ([]*IPRange, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	ipranges := make([]*IPRange, 0)
	scanner := bufio.NewScanner(f)
	lineno := 1
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		//comment start with '#'
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}
		ss := strings.Split(line, "-")
		if len(ss) != 2 {
			return nil, fmt.Errorf("Invalid line:%d in IP Range file:%s", lineno, file)
		}
		iprange, err := parseIPRange(ss[0], ss[1])
		if nil != err {
			return nil, fmt.Errorf("Invalid line:%d in IP Range file:%s", lineno, file)
		}
		ipranges = append(ipranges, iprange)
		lineno = lineno + 1
	}
	if len(ipranges) > 5 {
		dest := make([]*IPRange, len(ipranges))
		perm := rand.Perm(len(ipranges))
		for i, v := range perm {
			dest[v] = ipranges[i]
		}
		ipranges = dest
	}
	return ipranges, nil
}
