package main

import (
   "154.pages.dev/google/play"
   "154.pages.dev/http"
   "flag"
   "fmt"
   "strings"
)

type flags struct {
   acquire bool
   delivery bool
   app play.Application
   code string
   device bool
   single bool
   skipssl bool
   platform play.Platform
   bulkdetails bool
}

func main() {
   var f flags
   flag.StringVar(&f.app.ID, "a", "", "application ID")
   flag.BoolVar(&f.acquire, "acquire", false, "acquire application")
   flag.StringVar(&f.app.AssetModule, "asset", "", "download application asset module")
   flag.BoolVar(&f.bulkdetails, "bulkdetails", false, "fetch application details using /fdfe/bulkDetails request")
   flag.BoolVar(&f.delivery, "download", false, "download application")
   flag.BoolVar(&f.device, "d", false, "checkin and sync device")
   flag.StringVar(&f.code, "o", "", func() string {
      var b strings.Builder
      b.WriteString("oauth_token from ")
      b.WriteString("accounts.google.com/embedded/setup/v2/android")
      return b.String()
   }())
   flag.BoolVar(&f.single, "s", false, "single APK")
   flag.BoolVar(&f.skipssl, "sv", false, "skip SSL certificate verification for proxying")
   flag.Uint64Var(&f.app.Version, "v", 0, "version code")
   flag.Var(&f.platform, "p", fmt.Sprint(play.Platforms))
   flag.StringVar(&f.app.Languages, "l", "en-US,fr-FR,de-DE,it-IT,es-ES", "languages to download, comma separated")
   flag.Parse()
   http.No_Location()
   http.Verbose()
   if f.skipssl {
      http.SkipSSL()
   }
   switch {
   case f.app.ID != "":
      switch {
      case f.acquire:
         err := f.do_acquire()
         if err != nil {
            panic(err)
         }
      case f.app.Version >= 1:
         if f.delivery {
            err := f.do_delivery()
            if err != nil {
               panic(err)
            }
         } else if f.app.AssetModule != "" {
            err := f.do_asset_delivery()
            if err != nil {
               panic(err)
            }
         } else {
            details, err := f.do_details()
            if err != nil {
               panic(err)
            }
            fmt.Println(details)
         }
      default:
         details, err := f.do_details()
         if err != nil {
            panic(err)
         }
         fmt.Println(details)
      }
   case f.code != "":
      err := f.do_auth()
      if err != nil {
         panic(err)
      }
   case f.device:
      err := f.do_device()
      if err != nil {
         panic(err)
      }
   default:
      flag.Usage()
   }
}
