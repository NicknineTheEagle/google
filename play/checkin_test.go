package play

import (
   "fmt"
   "testing"
   "time"
)

func Test_Checkin(t *testing.T) {
   for _, platform := range Platforms {
      fmt.Println(platform)
      var check Checkin
      Phone.Platform = platform
      err := check.Checkin(Phone)
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
