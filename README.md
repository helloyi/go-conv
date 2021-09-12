# conv

Easy convert from one type to another in Go.

**NOTE: unstable**



## Why

If you are working with interfaces to handle things like dynamic content you’ll need an easy way to convert an interface into a given type. This is the library for you.

If you are taking in data from YAML, TOML or JSON or other formats which lack full types, then `conv` is the library for you.



## Example

Parse configuration with `viper` and `conv`.

```go
package main

import (
	"bytes"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"time"

	"github.com/helloyi/go-conv"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigType("toml")

	cfgData := []byte(`
Num = 21
String = "String"
Duration = "100ms"
Time = "Fri Nov 1 19:13:55 +0800 CST 2019"
ByteSize = "100MB"
IP = "8.8.8.8"
MAC = "02:00:5e:10:00:00:00:01"
URL = "http://host/path?param=x"
Regexp = "[0-9]+"
Mail = "name <user@mail.com>"
`)

	if err := viper.ReadConfig(bytes.NewBuffer(cfgData)); err != nil {
		panic(err)
	}

	rawCfg := viper.AllSettings()

	var cfg struct {
		Num      int
		String   string
		Duration time.Duration
		Time     *time.Time
		ByteSize conv.ByteSize
		IP       net.IP
		MAC      net.HardwareAddr
		URL      *url.URL
		Regexp   *regexp.Regexp
		Mail     *mail.Address
	}

	if err := conv.To(rawCfg, &cfg); err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", cfg)
	// Output: {21 String 100ms 2019-11-01 19:13:55 +0800 CST 104857600 8.8.8.8 02:00:5e:10:00:00:00:01 http://host/path?param=x [0-9]+ "name" <user@mail.com>}
}
```

