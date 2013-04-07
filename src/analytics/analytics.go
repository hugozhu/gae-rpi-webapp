package analytics

import (
	"log"
	"sync"
	"time"
)

var pv_lock = sync.Mutex{}

type Analytics struct {
	sync.Mutex
	pv *PV
	uv *UV
}

type PV struct {
	slots        []int
	base         time.Time //the first slot's timestamp
	interval     int       //the duration of slot
	num_of_slots int
	offset       int //indicate previous slot we were adding up
	all          []int64
}

type UV struct {
	all          map[string]int64
	interval     int //in seconds
	expiration   int //in seconds
	slots        []int
	num_of_slots int
}

func NewAnalytics(interval int, num_of_slots int) *Analytics {
	a := &Analytics{
		uv: NewUV(interval, num_of_slots),
		pv: NewPV(interval, num_of_slots),
	}
	return a
}

func (a *Analytics) AddOne(zcookie string, t time.Time) {
	a.Lock()
	defer a.Unlock()

	a.pv.AddOne(t)
	a.uv.AddOne(zcookie, t)
}

func (a *Analytics) Sum() (uv int, uv_slots []int, pv int, pv_slots []int) {
	a.Lock()
	defer a.Unlock()

	uv, uv_slots = a.uv.Sum()
	pv, pv_slots = a.pv.Sum()
	return
}

func NewUV(interval int, num_of_slots int) *UV {
	uv := &UV{
		interval:     interval,
		all:          map[string]int64{},
		expiration:   interval * num_of_slots,
		num_of_slots: num_of_slots,
	}
	return uv
}

func (uv *UV) AddOne(zcookie string, timestamp time.Time) {
	uv.all[zcookie] = timestamp.Unix()
}

func (uv *UV) Sum() (int, []int) {
	now := time.Now().Unix()
	slots := make([]int, uv.num_of_slots)
	for k, v := range uv.all {
		delta := int(now - v)
		if delta > uv.expiration {
			delete(uv.all, k)
			continue
		}
		p := uv.num_of_slots - 1 - delta/uv.interval
		if p < 0 {
			p = 0
		}
		slots[p]++
	}
	return len(uv.all), slots
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
	pv.all = []int64{}
	return pv
}

func (pv *PV) AddOne(t time.Time) {
	pv.all = append(pv.all, t.Unix())
}

func (pv *PV) Sum() (int, []int) {
	now := time.Now().Unix()
	threshold := pv.interval * pv.num_of_slots
	slots := make([]int, pv.num_of_slots)
	total := 0
	for i, v := range pv.all {
		if int(now-v) <= threshold {
			pv.all = pv.all[i:]
			total = len(pv.all)
			for i = 0; i < len(pv.all); i++ {
				p := pv.num_of_slots - 1 - int(now-pv.all[i])/pv.interval
				if p < 0 {
					p = 0
				}
				slots[p]++
			}
			break
		}
	}
	return total, slots
}

func (pv *PV) clear(from int, end int) {
	for i := from; i < end; i++ {
		pv.slots[i] = 0
	}
}

func (pv *PV) _add(timestamp time.Time, count int) {
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

func (pv *PV) _sum() (int, []int) {
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
