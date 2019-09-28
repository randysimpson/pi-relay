# Pi-Relay

This project will open a port `8080` to http traffic for a REST API interface to control
GPIO outputs for items hooked up to Raspberry Pi 3.

## Installation

The binaries can be downloaded by running:

```sh
wget https://github.com/randysimpson/pi-relay/releases/download/v1.0/pi-relay
chmod +x py-relay
```

To run the application:

```sh
./py-relay
```

## Operation

To use the API from the machine that is running the application:

```sh
curl localhost:8080/api/v1/gpio
```

### Logs

The logs are located at `app.log`.

## API

* `/api/v1/gpio` - Get

  Method used to get the GPIO outputs and states.  The `end_date` variable is only present when the state will be changing.  Therefore in the following example gpio with id `1ab2d` will be changing state to false in 5 minutes 30 seconds, and gpio with id `3vb7f` is in false state.

  ###### JSON response example:

  ```json
  {
    "status": "Success",
    "date": "2019-09-27T15:18:13-06:00Z",
    "result": [
      {
        "id": "1ab2d",
        "gpio": 4,
        "state": true,
        "end_date": "2019-09-27T15:23:43-06:00Z"
      },
      {
        "id": "3vb7f",
        "gpio": 17,
        "state": false
      }
    ]
  }
  ```

* `/api/v1/gpio` - Post

    Method used to create a GPIO on the device and set the initial state.

    The body can consist of 1 object or it can be an array of objects.

    ###### JSON body example:

    ```json
    {
      "id": "1vv3s",
      "gpio": 17,
      "state": false
    }
    ```

    ###### JSON response example:

    ```json
    {
      "status": "Created",
      "date": "2019-09-27T15:18:13-06:00Z",
      "result": {
        "id": "1vv3s",
        "gpio": 17,
        "state": false
      }
    }
    ```

* `/api/v1/gpio/{id}` - Get

    Method used to return details on a specific GPIO on the device.  The `end_date` variable is only present when the state will be changing.

    ###### JSON response example:

    ```json
    {
      "status": "Success",
      "date": "2019-09-27T15:18:13-06:00Z",
      "result": {
        "id": "1ab2d",
        "gpio": 4,
        "state": true,
        "end_date": "2019-09-27T15:23:43-06:00Z"
      }
    }
    ```

* `/api/v1/gpio/{id}` - Put

  Method used to modify the GPIO or to turn on/off a GPIO output on the device.
  The optional duration object can be in `minutes`/`seconds`/`milliseconds`/ or
  `end_date`.  If the duration object is omitted then the gpio will remain in the state until another put method is called to change it.

  > The GPIO location cannot be modified through this put method, it will only allow for changing state.

  ###### JSON body example:

  ```json
  {
    "id": "1ab2d",
    "gpio": 4,
    "state": true,
    "duration": {
      "minutes": 5,
      "seconds": 30
    }
  }
  ```

  ###### JSON response example:

  ```json
  {
    "status": "Updated",
    "date": "2019-09-27T15:18:13-06:00Z",
    "result": {
      "id": "1ab2d",
      "gpio": 4,
      "state": true,
      "end_date": "2019-09-27T15:23:43-06:00Z"
    }
  }
  ```

* `/api/v1/gpio/{id}` - Delete

    Method used to delete GPIO on the device.

    ###### JSON response example:

    ```json
    {
      "status": "Deleted",
      "date": "2019-09-27T15:18:13-06:00Z",
      "result": {
        "id": "1ab2d",
        "gpio": 4,
        "state": true
      }
    }
    ```

* `/api/v1/state` - Put

  Method used to modify the GPIO or to turn on/off a GPIO output on the device.
  The optional duration object can be in `minutes`/`seconds`/`milliseconds`/ or
  `end_date`.  If the duration object is omitted then the gpio will remain in the state until another put method is called to change it.

  > The GPIO location cannot be modified through this put method, it will only allow for changing state.

  The body can consist of 1 object or it can be an array of objects.

  ###### JSON body example:

  ```json
  [
    {
      "id": "1ab2d",
      "gpio": 4,
      "state": true,
      "duration": {
        "minutes": 5,
        "seconds": 30
      }
    }, {
      "id": "1vv3s",
      "gpio": 17,
      "state": true
    }
  ]
  ```

  ###### JSON response example:

  ```json
  {
    "status": "Updated",
    "date": "2019-09-27T15:18:13-06:00Z",
    "result": [
      {
        "id": "1ab2d",
        "gpio": 4,
        "state": true,
        "end_date": "2019-09-27T15:23:43-06:00Z"
      }, {
        "id": "1vv3s",
        "gpio": 17,
        "state": true
      }
    ]
  }
  ```

### Extra

#### Dependencies

For rpio

```sh
go get github.com/stianeikeland/go-rpio
```

For API

```sh
go get -u github.com/gorilla/mux
```

For error logging

```sh
go get github.com/pkg/errors
```

#### Build

```sh
env GOOS=linux GOARCH=arm GOARM=5 go build
```

#### Run

```sh
./pi-relay
```

## License

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
