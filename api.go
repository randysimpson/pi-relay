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
  "log"
  "time"
  "bytes"
  "net/http"
  "encoding/json"
  "io/ioutil"
  "github.com/pkg/errors"
  "github.com/gorilla/mux"
)

type putGpio struct {
  ID  string `json:"id"`
  Gpio int `json:"gpio"`
  State bool `json:"state"`
  Duration duration `json:"duration,omitempty"`
}

type duration struct {
  Minutes int `json:"minutes,omitempty"`
  Seconds int `json:"seconds,omitempty"`
  Milliseconds int `json:"milliseconds,omitempty"`
  EndDate string `json:"end_date,omitempty"`
}

func (p* putGpio) GetDurationMS() (int64, error) {
  var rtnMS int64
  rtnMS = 0
  if p.Duration.EndDate != "" {
    endDateStr := p.Duration.EndDate
    if endDateStr[len(endDateStr)-1:] == "Z" {
      endDateStr = endDateStr[:len(endDateStr)-1]
    }
    end, err := time.Parse(time.RFC3339, endDateStr)
    if err != nil {
      return -1, errors.Wrap(err, fmt.Sprintf("unable to parse end date %s", p.Duration.EndDate))
    }
    rtnMS = time.Until(end).Milliseconds()
  } else {
      rtnMS += int64(p.Duration.Minutes * 60000)
      rtnMS += int64(p.Duration.Seconds * 1000)
      rtnMS += int64(p.Duration.Milliseconds)
  }
  return rtnMS, nil
}

type allGpio []Gpio

var gpios = allGpio{
}

func isJsonArray(data []byte) bool {
  // Get slice of data with optional leading whitespace removed.
  // See RFC 7159, Section 2 for the definition of JSON whitespace.
  x := bytes.TrimLeft(data, " \t\r\n")

  return len(x) > 0 && x[0] == '['
}

func homeLink(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Return API specification (Future).")
}

func getAllGpio(w http.ResponseWriter, r *http.Request) {
  rtnJson, _ := json.Marshal(gpios)

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)

  logMessage(fmt.Sprintf("%sZ - %s - %s%s - 200 - %s\n", time.Now().Format(time.RFC3339), r.Method, r.Host, r.URL.Path, r.RemoteAddr))
  w.Write([]byte(fmt.Sprintf(`{ "status":"Success", "date":"%sZ", "result":%s }`, time.Now().Format(time.RFC3339), string(rtnJson))))
}

func setState(g *Gpio, updatedGpio putGpio) error {
  durationMS, err := updatedGpio.GetDurationMS()
  if err != nil {
    return errors.Wrap(err, fmt.Sprint("unable to get duration"))
  }
  ptrDuration := &durationMS
  if durationMS <= 0 {
    ptrDuration = nil
  }
  g.Toggle(updatedGpio.State, ptrDuration)
  return nil
}

func createGpio(w http.ResponseWriter, r *http.Request) {
  reqBody, err := ioutil.ReadAll(r.Body)
  if err != nil {
    badRequest(w, r, err)
    return
  }

  var rtnJson []byte
  isArray := isJsonArray(reqBody)
  if isArray {
    var newPutGpioList []putGpio
    json.Unmarshal(reqBody, &newPutGpioList)

    for i := 0; i < len(newPutGpioList); i++ {
      var newGpio = Gpio{
        ID: newPutGpioList[i].ID,
        Gpio: newPutGpioList[i].Gpio,
        State:false,
      }

      //set initial state:
      err := setState(&newGpio, newPutGpioList[i])
      if err != nil {
        logMessage(fmt.Sprintf("create error %+v", err))
      }

      gpios = append(gpios, newGpio)
    }

    rtnJson, _ = json.Marshal(newPutGpioList)
  } else {
    var newPutGpio putGpio
    json.Unmarshal(reqBody, &newPutGpio)

    var newGpio = Gpio{
      ID:newPutGpio.ID,
      Gpio:newPutGpio.Gpio,
      State:false,
    }
    gpios = append(gpios, newGpio)

    //set initial state
    err := setState(&gpios[len(gpios) - 1], newPutGpio)
    if err != nil {
      logMessage(fmt.Sprintf("create error %+v", err))
    }

    rtnJson, _ = json.Marshal(newGpio)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusCreated)

  logMessage(fmt.Sprintf("%sZ - %s - %s%s - 201 - %s\n", time.Now().Format(time.RFC3339), r.Method, r.Host, r.URL.Path, r.RemoteAddr))
  w.Write([]byte(fmt.Sprintf(`{ "status":"Created", "date":"%sZ", "result":%s }`, time.Now().Format(time.RFC3339), string(rtnJson))))
}

func getGpio(w http.ResponseWriter, r *http.Request) {
  id := mux.Vars(r)["id"]

  for _, gpio := range gpios {
    if gpio.ID == id {
      w.Header().Set("Content-Type", "application/json; charset=UTF-8")
      w.WriteHeader(http.StatusOK)

      logMessage(fmt.Sprintf("%sZ - %s - %s%s - 200 - %s\n", time.Now().Format(time.RFC3339), r.Method, r.Host, r.URL.Path, r.RemoteAddr))
      rtnJson, _ := json.Marshal(gpio)
      w.Write([]byte(fmt.Sprintf(`{ "status":"Success", "date":"%sZ", "result":%s }`, time.Now().Format(time.RFC3339), string(rtnJson))))
      return
    }
  }

  //it was not found.
  notFound(w, r)
}

func updateGpio(w http.ResponseWriter, r *http.Request) {
  id := mux.Vars(r)["id"]
  var updatedGpio putGpio

  reqBody, err := ioutil.ReadAll(r.Body)
  if err != nil {
    badRequest(w, r, err)
    return
  }

  json.Unmarshal(reqBody, &updatedGpio)

  for i := 0; i < len(gpios); i++ {
    if gpios[i].ID == id {
      err := setState(&gpios[i], updatedGpio)
      if err != nil {
        logMessage(fmt.Sprintf("state error %+v", err))
      }

      rtnJson, _ := json.Marshal(gpios[i])

      w.Header().Set("Content-Type", "application/json; charset=UTF-8")
      w.WriteHeader(http.StatusOK)

      logMessage(fmt.Sprintf("%sZ - %s - %s%s - 200 - %s\n", time.Now().Format(time.RFC3339), r.Method, r.Host, r.URL.Path, r.RemoteAddr))
      w.Write([]byte(fmt.Sprintf(`{ "status":"Updated", "date":"%sZ", "result":%s }`, time.Now().Format(time.RFC3339), rtnJson)))
      return
    }
  }

  //it was not found.
  notFound(w, r)
}

func deleteGpio(w http.ResponseWriter, r *http.Request) {
  id := mux.Vars(r)["id"]

  for i, gpio := range gpios {
    if gpio.ID == id {
      gpios = append(gpios[:i], gpios[i+1:]...)

      w.Header().Set("Content-Type", "application/json; charset=UTF-8")
      w.WriteHeader(http.StatusOK)

      logMessage(fmt.Sprintf("%sZ - %s - %s%s - 200 - %s\n", time.Now().Format(time.RFC3339), r.Method, r.Host, r.URL.Path, r.RemoteAddr))
      rtnJson, _ := json.Marshal(gpio)
      w.Write([]byte(fmt.Sprintf(`{ "status":"Deleted", "date":"%sZ", "result":%s }`, time.Now().Format(time.RFC3339), rtnJson)))
      return
    }
  }

  //it was not found.
  notFound(w, r)
}

func changeState(w http.ResponseWriter, r *http.Request) {
  reqBody, err := ioutil.ReadAll(r.Body)
  if err != nil {
    badRequest(w, r, err)
    return
  }

  rtnJson := ""
  var updatedGpioList []putGpio
  isArray := isJsonArray(reqBody)
  if isArray {
    json.Unmarshal(reqBody, &updatedGpioList)
  } else {
    updatedGpioList = []putGpio{

    }
    var updatedGpio putGpio
    json.Unmarshal(reqBody, &updatedGpio)
    updatedGpioList = append(updatedGpioList, updatedGpio)
  }

  updateCount := 0
  for i := 0; i < len(gpios); i++ {
    for j := 0; j < len(updatedGpioList); j++ {
      if gpios[i].ID == updatedGpioList[j].ID {
        err := setState(&gpios[i], updatedGpioList[j])
        if err != nil {
          logMessage(fmt.Sprintf("state error %+v", err))
        }

        itemJson, _ := json.Marshal(gpios[i])
        rtnJson = rtnJson + "," + string(itemJson)
        updateCount ++
      }
    }
  }

  if updateCount == 0 {
    //it was not found.
    notFound(w, r)
    return
  } else if updateCount > 1 {
    rtnJson = "[" + rtnJson[1:] + "]"
  } else if updateCount == 1 {
    rtnJson = rtnJson[1:]
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)

  logMessage(fmt.Sprintf("%sZ - %s - %s%s - 200 - %s\n", time.Now().Format(time.RFC3339), r.Method, r.Host, r.URL.Path, r.RemoteAddr))
  w.Write([]byte(fmt.Sprintf(`{ "status":"Updated", "date":"%sZ", "result":%s }`, time.Now().Format(time.RFC3339), rtnJson)))
}

func badRequest(w http.ResponseWriter, r *http.Request, err error) {
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusBadRequest)
  logMessage(fmt.Sprintf("%sZ - %s - %s%s - 404 - %s - %+v\n", time.Now().Format(time.RFC3339), r.Method, r.Host, r.URL.Path, r.RemoteAddr, err))
  w.Write([]byte(fmt.Sprintf(`{ "status":"Error on Request Body", "date":"%sZ" }`, time.Now().Format(time.RFC3339))))
}

func notFound(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusNotFound)
  logMessage(fmt.Sprintf("%sZ - %s - %s%s - 404 - %s\n", time.Now().Format(time.RFC3339), r.Method, r.Host, r.URL.Path, r.RemoteAddr))
  w.Write([]byte(fmt.Sprintf(`{ "status":"Not Found", "date":"%sZ" }`, time.Now().Format(time.RFC3339))))
}

func logMessage(message string) {
  fmt.Print(message)
  log.Print(message)
}

func HandleRequest() {
  router := mux.NewRouter().StrictSlash(true)
  router.HandleFunc("/", homeLink)
  router.HandleFunc("/api/v1/gpio", createGpio).Methods("POST")
  router.HandleFunc("/api/v1/gpio", getAllGpio).Methods("GET")
  router.HandleFunc("/api/v1/gpio/{id}", getGpio).Methods("GET")
  router.HandleFunc("/api/v1/gpio/{id}", updateGpio).Methods("PUT")
  router.HandleFunc("/api/v1/gpio/{id}", deleteGpio).Methods("DELETE")
  router.HandleFunc("/api/v1/state", changeState).Methods("PUT")
  http.ListenAndServe(":8080", router)
}
