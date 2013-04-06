package analytics

import (
	"log"
	"sync"
	"time"
)

var pv_lock = sync.Mutex{}

type PV struct {
	sync.Mutex
	slots        []int
	base         time.Time //the first slot's timestamp
	interval     int       //the duration of slot
	num_of_slots int
	offset       int //indicate previous slot we were adding up
}

type UV struct {
	sync.Mutex
	all        map[string]int64
	interval   int //in seconds
	expiration int //in seconds
	slots      []int
}

func NewUV(interval int, num_of_slots int) *UV {
	uv := &UV{
		interval:   interval,
		all:        map[string]int64{},
		expiration: interval * num_of_slots,
		slots:      make([]int, num_of_slots),
	}
	return uv
}

func (uv *UV) AddOne(zcookie string, timestamp time.Time) {
	uv.Lock()
	defer uv.Unlock()
	uv.all[zcookie] = timestamp.Unix()
}

func (uv *UV) Sum() (int, []int) {
	uv.Lock()
	defer uv.Unlock()
	now := time.Now().Unix()
	num_of_slots := len(uv.slots)
	for i, _ := range uv.slots {
		uv.slots[i] = 0
	}
	for k, v := range uv.all {
		delta := int(now - v)
		if delta > uv.expiration {
			delete(uv.all, k)
			continue
		}
		p := num_of_slots - 1 - delta/uv.interval
		if p < 0 {
			p = 0
		}
		uv.slots[p]++
	}
	return len(uv.all), uv.slots
}

func NewPV(interval int, num_of_slots int) *PV {
	slots := make([]int, num_of_slots)
	pv := &PV{
		slots:        slots,
		base:         time.Now(),
		interval:     interval,
		num_of_slots: num_of_slots,
	}
	pv.clear(0, len(pv.slots))
	return pv
}

func (pv *PV) clear(from int, end int) {
	for i := from; i < end; i++ {
		pv.slots[i] = 0
	}
}

func (pv *PV) AddOne(t time.Time) {
	pv.Add(t, 1)
}

func (pv *PV) Add(timestamp time.Time, count int) {
	pv.Lock()
	defer pv.Unlock()

	now := timestamp
	delta := int(now.Unix()-pv.base.Unix()) / pv.interval
	index := delta % pv.num_of_slots

	if delta >= pv.num_of_slots*2 {
		pv.base = now
		pv.clear(0, len(pv.slots))
		pv.slots[index] = count
		pv.offset = index
	} else if delta >= pv.num_of_slots {
		pv.base = now
		pv.clear(0, delta-pv.num_of_slots)
		pv.slots[index] = count
		pv.offset = index
	} else {
		if index == pv.offset {
			pv.slots[index]++
		} else {
			pv.clear(pv.offset+1, index)
			pv.slots[index] = count
			pv.offset = index
		}
	}
}

func (pv *PV) Sum() (int, []int) {
	pv.Lock()
	defer pv.Unlock()
	now := time.Now()
	sum := 0
	delta := int(now.Unix()-pv.base.Unix()) / pv.interval

	slots := make([]int, pv.num_of_slots)

	log.Println(delta)
	if delta >= pv.num_of_slots*2 {
		return 0, slots
	} else if delta >= pv.num_of_slots {
		for i := delta - pv.num_of_slots; i < len(pv.slots); i++ {
			sum = sum + pv.slots[i]
			p := delta - pv.num_of_slots - i
			slots[p] = pv.slots[i]
		}
	} else {
		for i := 0; i <= pv.offset; i++ {
			p := pv.num_of_slots - delta - i - 1
			if p < 0 {
				p = pv.num_of_slots - 1
			}
			slots[p] = pv.slots[i]
			sum = sum + pv.slots[i]
		}
	}
	return sum, slots
}
