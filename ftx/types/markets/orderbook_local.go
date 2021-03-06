package markets

import (
	"sort"
	"sync"
	"time"
)

type OrderBookLocal struct {
	ob map[float64]*OrderBookL2
	m  sync.Mutex
}

func (o *OrderBookLocal) GetOrderBook() (ob OrderBook) {
	for _, v := range o.ob {
		switch v.Side {
		case "bid":
			ob.Bids = append(ob.Bids, Item{
				Price:  v.Price,
				Amount: v.Size,
			})
		case "ask":
			ob.Asks = append(ob.Asks, Item{
				Price:  v.Price,
				Amount: v.Size,
			})
		}
	}

	sort.Slice(ob.Bids, func(i, j int) bool {
		return ob.Bids[i].Price > ob.Bids[j].Price
	})

	sort.Slice(ob.Asks, func(i, j int) bool {
		return ob.Asks[i].Price < ob.Asks[j].Price
	})

	ob.Timestamp = time.Now()

	return
}

func NewOrderBookLocal() *OrderBookLocal {
	o := &OrderBookLocal{
		ob: make(map[float64]*OrderBookL2),
	}
	return o
}

func (o *OrderBookLocal) LoadSnapshot(newOrderbook *OrderBookRaw) error {
	o.m.Lock()
	defer o.m.Unlock()

	for _, ask := range newOrderbook.Asks {
		o.ob[ask[0]] = &OrderBookL2{Price: ask[0], Size: ask[1], Side: "ask"}
	}
	for _, bid := range newOrderbook.Bids {
		o.ob[bid[0]] = &OrderBookL2{Price: bid[0], Size: bid[1], Side: "bid"}
	}

	return nil
}

func (o *OrderBookLocal) Update(deltaOrderbook *OrderBookRaw) {
	o.m.Lock()
	defer o.m.Unlock()

	for _, item := range deltaOrderbook.Asks {
		if _, ok := o.ob[item[0]]; ok {
			if item[1] == 0.0 {
				// delete case
				delete(o.ob, item[0])
			} else {
				// update case
				o.ob[item[0]].Size = item[1]
			}
		} else {
			// insert case
			o.ob[item[0]] = &OrderBookL2{Price: item[0], Size: item[1], Side: "ask"}
		}
	}

	for _, item := range deltaOrderbook.Bids {
		if _, ok := o.ob[item[0]]; ok {
			if item[1] == 0.0 {
				// delete case
				delete(o.ob, item[0])
			} else {
				// update case
				o.ob[item[0]].Size = item[1]
			}
		} else {
			// insert case
			o.ob[item[0]] = &OrderBookL2{Price: item[0], Size: item[1], Side: "bid"}
		}
	}
}
