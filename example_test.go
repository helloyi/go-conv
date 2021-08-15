package conv_test

import (
	"bytes"
	"fmt"
	"time"

	"github.com/helloyi/go-conv"
	"github.com/spf13/viper"
)

func ExampleTo() {
	viper.SetConfigType("toml")

	cfgData := []byte(`
Num = 21
String = "String"
Duration = "100ms"
Time = "Fri Nov 1 19:13:55 +0800 CST 2019"
ByteSize = "100MB"
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
	}

	if err := conv.To(rawCfg, &cfg); err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", cfg)
	// Output: {21 String 100ms 2019-11-01 19:13:55 +0800 CST 104857600}
}
