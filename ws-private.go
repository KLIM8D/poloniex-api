package poloniex

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	turnpike "gopkg.in/beatgammit/turnpike.v2"
)

type (
	subscription struct {
		Command string `json:"command"`
		Channel string `json:"channel"`
	}
)

func (p *Poloniex) sendWSMessage(msg interface{}) error {

	msgs, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "marshalling WSmessage failed")
	}
	log.Println(string(msgs))

	err = p.ws.WriteMessage(websocket.TextMessage, msgs)
	if err != nil {
		return errors.Wrap(err, "sending WSmessage failed")
	}
	return nil
}

//makeTickerHandler takes a WS Order or Trade and send it over the channel sepcified by the user
func (p *Poloniex) messageHandler(ch chan WSTicker) turnpike.EventHandler {
	return func(p []interface{}, n map[string]interface{}) {
		t := WSTicker{
			Pair:          p[0].(string),
			Last:          toFloat(p[1]),
			Ask:           toFloat(p[2]),
			Bid:           toFloat(p[3]),
			PercentChange: toFloat(p[4]) * 100.0,
			BaseVolume:    toFloat(p[5]),
			QuoteVolume:   toFloat(p[6]),
			IsFrozen:      toFloat(p[7]) != 0.0,
			DailyHigh:     toFloat(p[8]),
			DailyLow:      toFloat(p[9]),
		}
		ch <- t
	}
}
