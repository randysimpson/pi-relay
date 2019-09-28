/*
pi-relay

MIT License

Copyright (©) 2019 - Randall Simpson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package main

import (
  "fmt"
  "time"
  "github.com/pkg/errors"
  "github.com/stianeikeland/go-rpio"
)

type Gpio struct {
  ID  string `json:"id"`
  Gpio int `json:"gpio"`
  State bool `json:"state"`
  EndDate string `json:"end_date,omitempty"`
}

func (g *Gpio) Toggle(state bool, duration *int64) error {
  err := rpio.Open()
  if err != nil {
    return errors.Wrap(err, fmt.Sprint("unable to open gpio %s", g.ID))
  }

  defer rpio.Close()

  pin := rpio.Pin(g.Gpio)
  pin.Output()

  if state == true {
    pin.Low()
  } else {
    pin.High()
  }
  g.State = state

  if duration != nil {
    end := time.Now().Add(time.Millisecond * time.Duration(*duration))
    g.EndDate = fmt.Sprintf("%sZ", end.Format(time.RFC3339))
    go delayToggle(!state, g, duration)
  }
  return nil
}

func delayToggle(state bool, g *Gpio, duration *int64) {
  time.Sleep(time.Millisecond * time.Duration(*duration))
  //check to see if the toggle is still valid.
  if g.EndDate != "" {
    g.Toggle(state, nil)
    g.EndDate = ""
  }
}
