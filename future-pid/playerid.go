package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type playerId [4]uint64

func (p playerId) MarshalJSON() ([]byte, error) {
	o := make([]string, len(p), len(p))
	for i, s := range p {
		o[i] = strconv.FormatUint(s, 36)
	}
	return json.Marshal(strings.Join(o, "-"))
}

func (p *playerId) UnmarshalJSON(s []byte) error {
	str := ""
	err := json.Unmarshal(s, &str)
	if err != nil {
		return err
	}
	strs := strings.Split(str, "-")
	if len(strs) != len(p) {
		return fmt.Errorf("unable to parse %v as a playerId, expected %d segments, got %d", str, len(p), len(strs))
	}
	for i, seg := range strs {
		segi, err := strconv.ParseUint(seg, 36, 64)
		if err != nil {
			return fmt.Errorf("unable to parse %v as a playerId, segment %v, reason: %v", str, seg, err.Error())
		}
		p[i] = segi
	}
	return nil
}

func main() {
	thing1 := playerId{}
	for i := range thing1 {
		thing1[i] = rand.Uint64()
	}
	fmt.Printf("%+v\n", thing1)
	js, err := json.Marshal(thing1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(js))
	pid := &playerId{}
	err = json.Unmarshal(js, pid)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%+v\n", *pid)
}
