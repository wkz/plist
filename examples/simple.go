package main

import (
	"fmt"
	"time"

	"github.com/wkz/plist"
)

type Person struct {
	Name           string
	Birth          time.Time `plist:"BirthDate"`
	Height         float32
	DriversLicense bool
	PhoneNumbers   map[string]int
	FingerPrints   [5][]byte
	TrustedByMe    bool `plist:"-"`
}

func main() {
	jd := Person{
		Name: "John Doe",
		Birth: time.Date(1985, time.June, 20, 0, 0, 0, 0, time.Local),
		Height: 180.5,
		DriversLicense: true,
		PhoneNumbers: map[string]int{"home": 5550190, "work": 5550001},
		FingerPrints: [...][]byte{
			[]byte("7hum6"), 
			[]byte("1ndex"),
			[]byte("l0ng"),
			[]byte("r1ng"), 
			[]byte("p1nky"),
		},
		TrustedByMe: false,
	}

	plist, err := plist.MarshalIndent(&jd, "", "\t")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Println(string(plist))
}