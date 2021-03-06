// Copyright 2013 Mikio Hara. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.

package ipaddr

import (
	"bytes"
	"net"
	"sort"
)

type byAddrFamily []Prefix

func (ps byAddrFamily) newIPv4Prefixes() []Prefix {
	nps := make([]Prefix, 0, len(ps))
	for _, p := range ps {
		if p.IP.To4() != nil {
			np := clonePrefix(&p)
			nps = append(nps, *np)
		}
	}
	return nps
}

func (ps byAddrFamily) newIPv6Prefixes() []Prefix {
	nps := make([]Prefix, 0, len(ps))
	for _, p := range ps {
		if p.IP.To16() != nil && p.IP.To4() == nil {
			np := clonePrefix(&p)
			nps = append(nps, *np)
		}
	}
	return nps
}

type byAscending []Prefix

func (ps byAscending) Len() int           { return len(ps) }
func (ps byAscending) Less(i, j int) bool { return compareAscending(&ps[i], &ps[j]) < 0 }
func (ps byAscending) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }

func compareAscending(a, b *Prefix) int {
	if n := bytes.Compare(a.IP, b.IP); n != 0 {
		return n
	}
	if n := bytes.Compare(a.Mask, b.Mask); n != 0 {
		return n
	}
	return 0
}

type byDescending []Prefix

func (ps byDescending) Len() int           { return len(ps) }
func (ps byDescending) Less(i, j int) bool { return compareDescending(&ps[i], &ps[j]) >= 0 }
func (ps byDescending) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }

func compareDescending(a, b *Prefix) int {
	if n := bytes.Compare(a.Mask, b.Mask); n != 0 {
		return n
	}
	if n := bytes.Compare(a.IP, b.IP); n != 0 {
		return n
	}
	return 0
}

type sortDir int

const (
	sortAscending sortDir = iota
	sortDescending
)

func newSortedPrefixes(ps []Prefix, dir sortDir, strict bool) []Prefix {
	if len(ps) == 0 {
		return nil
	}
	if strict {
		if ps[0].IP.To4() != nil {
			ps = byAddrFamily(ps).newIPv4Prefixes()
		}
		if ps[0].IP.To16() != nil && ps[0].IP.To4() == nil {
			ps = byAddrFamily(ps).newIPv6Prefixes()
		}
		if dir == sortAscending {
			sort.Sort(byAscending(ps))
		} else {
			sort.Sort(byDescending(ps))
		}
	} else {
		nps := make([]Prefix, 0, len(ps))
		for i := range ps {
			np := clonePrefix(&ps[i])
			nps = append(nps, *np)
		}
		if dir == sortAscending {
			sort.Sort(byAscending(nps))
		} else {
			sort.Sort(byDescending(nps))
		}
		ps = nps
	}
	nps := ps[:0]
	var p *Prefix
	for i := range ps {
		if p == nil {
			nps = append(nps, ps[i])
		} else if !p.Equal(&ps[i]) {
			nps = append(nps, ps[i])
		}
		p = &ps[i]
	}
	return nps
}

func clonePrefix(s *Prefix) *Prefix {
	d := &Prefix{IPNet: net.IPNet{IP: make(net.IP, net.IPv6len), Mask: make(net.IPMask, len(s.Mask))}}
	copy(d.IP, s.IP.To16())
	copy(d.Mask, s.Mask)
	return d
}
