package analytics

import (
	"log"
	"math/rand"
	"testing"
	"time"
)

func _TestUV(t *testing.T) {
	uv := NewUV(30)
	count := 1
	visitors := []string{"a", "b", "c", "d", "e", "f"}
	go func() {
		for {
			log.Println(count, uv.Sum())
			count++
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		for i := 0; i < 15; i++ {
			cookie := visitors[rand.Intn(len(visitors))]
			uv.AddOne(cookie, time.Now())
			time.Sleep(1 * time.Second)
		}
	}()
	time.Sleep(2 * time.Minute)
	t.Log(uv)
}

var logs []int64

func correct_pv(threshold int64) int {
	count := 0
	now := time.Now().Unix()
	for _, t := range logs {
		if now-t <= threshold {
			count++
		}
	}
	return count
}

func TestPV(t *testing.T) {
	visitors := []string{"a", "b", "c", "d", "e", "f"}

	pv := NewPV(2, 10) //count pv in 20 seconds
	uv := NewUV(2 * 10)
	count := 1
	go func() {
		for {
			log.Println(count, "[", uv.Sum(), ",", pv.Sum(), ",", correct_pv(20), "]", pv.slots, uv.all, pv.base.Unix(), time.Now().Unix())
			count++
			if pv.Sum() > 20 {
				t.Error("Impossible to exceed 20 PVs")
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for i := 0; i < 15; i++ {
			now := time.Now()
			pv.AddOne(now)
			zcookie := visitors[rand.Intn(len(visitors))]
			uv.AddOne(zcookie, now)
			logs = append(logs, now.Unix())
			time.Sleep(1 * time.Second)
		}
		time.Sleep(60 * time.Second)
		for i := 0; i < 10; i++ {
			now := time.Now()
			pv.AddOne(now)
			zcookie := visitors[rand.Intn(len(visitors))]
			uv.AddOne(zcookie, now)
			logs = append(logs, now.Unix())
			time.Sleep(20 * time.Second)
		}
		time.Sleep(41 * time.Second)
		for i := 0; i < 15; i++ {
			now := time.Now()
			pv.AddOne(now)
			zcookie := visitors[rand.Intn(len(visitors))]
			uv.AddOne(zcookie, now)
			logs = append(logs, now.Unix())
			time.Sleep(10 * time.Second)
		}
	}()

	time.Sleep(2 * time.Minute)
	t.Log(pv)
}
