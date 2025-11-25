package http

import (
   "154.pages.dev/encoding/hex"
   "fmt"
   "crypto/tls"
   "net/http"
   "net/http/httputil"
)

func SkipSSL() func() {
   k := &http.DefaultClient.Transport
   v := *k
   tr := &http.Transport{
      TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
   }
   *k = tr
   return func() { *k = v }
}

func Location() func() {
   k := &http.DefaultClient.CheckRedirect
   v := *k
   *k = nil
   return func() { *k = v }
}

func No_Location() func() {
   k := &http.DefaultClient.CheckRedirect
   v := *k
   *k = func(*http.Request, []*http.Request) error {
      return http.ErrUseLastResponse
   }
   return func() { *k = v }
}

func Silent() func() {
   k := &http.DefaultClient.Transport
   v := *k
   *k = nil
   return func() { *k = v }
}

type transport func(*http.Request) (int, error)

func (t transport) RoundTrip(req *http.Request) (*http.Response, error) {
   _, err := t(req)
   if err != nil {
      return nil, err
   }
   return http.DefaultTransport.RoundTrip(req)
}

func Verbose() func() {
   k := &http.DefaultClient.Transport
   v := *k
   *k = transport(func(r *http.Request) (int, error) {
      return fmt.Println(r.Method, r.URL)
   })
   return func() { *k = v }
}

func Trace() func() {
   k := &http.DefaultClient.Transport
   v := *k
   *k = transport(func(r *http.Request) (int, error) {
      b, err := httputil.DumpRequest(r, true)
      if err != nil {
         return 0, err
      }
      if hex.Binary(b) {
         b = hex.Encode(b)
      }
      return fmt.Println(string(b))
   })
   return func() { *k = v }
}
