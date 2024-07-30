package play

import (
   "154.pages.dev/protobuf"
   "bytes"
   "compress/gzip"
   "encoding/base64"
   "net/http"
   "strconv"
   "time"
)

const (
   android_API = 30
   google_play_store = 82941300
   // developer.android.com/guide/topics/manifest/uses-feature-element#glEsVersion
   // the device actually uses 0x30000, but some apps require a higher version,
   // so we need to lie:
   gl_es_version = 0x30001
)

var Phone = Device{
   Texture: []string{
      // com.instagram.android
      "GL_OES_compressed_ETC1_RGB8_texture",
      // com.kakaogames.twodin
      "GL_KHR_texture_compression_astc_ldr",
      // com.sega.sonic4ep2thd
      "GL_EXT_texture_compression_s3tc",
      // com.sonicjump.sonicjump
      "GL_AMD_compressed_3DC_texture",
      "GL_AMD_compressed_ATC_texture",
      "GL_ATI_texture_compression_atitc",
      "GL_IMG_texture_compression_pvrtc",
      // net.fishlabs.GalaxyonFire2THD
      "GL_EXT_texture_compression_dxt1",
   },
   Library: []string{
      // com.amctve.amcfullepisodes
      "org.apache.http.legacy",
      // com.binance.dev
      "android.test.runner",
   },
   Feature: []string{
      // app.source.getcontact
      "android.hardware.location.gps",
      // br.com.rodrigokolb.realdrum
      "android.software.midi",
      // com.app.xt
      "android.hardware.camera.front",
      // com.cabify.rider
      "android.hardware.camera.flash",
      // com.clearchannel.iheartradio.controller
      "android.hardware.microphone",
      // com.google.android.apps.walletnfcrel
      "android.software.device_admin",
      // com.google.android.youtube
      "android.hardware.touchscreen",
      "android.hardware.wifi",
      // com.madhead.tos.zh
      "android.hardware.sensor.accelerometer",
      // com.pinterest
      "android.hardware.camera",
      "android.hardware.location",
      "android.hardware.screen.portrait",
      // com.supercell.brawlstars
      "android.hardware.touchscreen.multitouch",
      // com.mediocre.smashhit
      "android.hardware.touchscreen.multitouch.distinct",
      // com.sega.sonicjumpfever
      "android.hardware.touchscreen.multitouch.jazzhand",
      // com.sygic.aura
      "android.hardware.location.network",
      // com.xiaomi.smarthome
      "android.hardware.bluetooth",
      "android.hardware.bluetooth_le",
      "android.hardware.camera.autofocus",
      "android.hardware.usb.host",
      // kr.sira.metal
      "android.hardware.sensor.compass",
      // org.thoughtcrime.securesms
      "android.hardware.telephony",
      // org.videolan.vlc
      "android.hardware.screen.landscape",
   },
}

func authorization(r *http.Request, a Access_Token) {
   r.Header.Set("Authorization", "Bearer " + a.Values.Get("Auth"))
}

func compress(m protobuf.Message) (string, error) {
   var b bytes.Buffer
   w := gzip.NewWriter(&b)
   _, err := w.Write(m.Append(nil))
   if err != nil {
      return "", err
   }
   if err := w.Close(); err != nil {
      return "", err
   }
   return base64.URLEncoding.EncodeToString(b.Bytes()), nil
}

func user_agent(r *http.Request, single bool) {
   var b []byte
   // `sdk` is needed for `/fdfe/delivery`
   b = append(b, "Android-Finsky (sdk="...)
   // with `/fdfe/acquire`, requests will be rejected with certain apps, if the
   // device was created with too low a version here:
   b = strconv.AppendInt(b, android_API, 10)
   b = append(b, ",versionCode="...)
   // for multiple APKs just tell the truth. for single APK we have to lie.
   // below value is the last version that works.
   if single {
      b = strconv.AppendInt(b, 80919999, 10)
   } else {
      b = strconv.AppendInt(b, google_play_store, 10)
   }
   b = append(b, ')')
   r.Header.Set("User-Agent", string(b))
}

func x_dfe_device_id(r *http.Request, c Checkin) error {
   id, err := c.Device_ID()
   if err != nil {
      return err
   }
   r.Header.Set("X-DFE-Device-ID", strconv.FormatUint(id, 16))
   return nil
}

func x_ps_rh(r *http.Request, c Checkin) error {
   id, err := c.Device_ID()
   if err != nil {
      return err
   }
   token, err := func() (string, error) {
      var m protobuf.Message
      m.Add(3, func(m *protobuf.Message) {
         m.Add_String(1, strconv.FormatUint(id, 10))
         m.Add(2, func(m *protobuf.Message) {
            v := time.Now().UnixMicro()
            m.Add_String(1, strconv.FormatInt(v, 10))
         })
      })
      return compress(m)
   }()
   if err != nil {
      return err
   }
   ps_rh, err := func() (string, error) {
      var m protobuf.Message
      m.Add(1, func(m *protobuf.Message) {
         m.Add_String(1, token)
      })
      return compress(m)
   }()
   if err != nil {
      return err
   }
   r.Header.Set("X-PS-RH", ps_rh)
   return nil
}

type Application struct {
   ID string
   Version uint64
}

func (a Application) APK(config string) string {
   var b []byte
   b = append(b, a.ID...)
   b = append(b, '-')
   if config != "" {
      b = append(b, config...)
      b = append(b, '-')
   }
   b = strconv.AppendUint(b, a.Version, 10)
   b = append(b, ".apk"...)
   return string(b)
}

func (a Application) OBB(role uint64, version_code uint64) string {
   var b []byte
   if role >= 1 {
      b = append(b, "patch"...)
   } else {
      b = append(b, "main"...)
   }
   b = append(b, '.')
   b = strconv.AppendUint(b, version_code, 10)
   b = append(b, '.')
   b = append(b, a.ID...)
   b = append(b, ".obb"...)
   return string(b)
}

type Platform int

var Platforms = map[int]string{
   // com.google.android.youtube
   0: "x86",
   // com.miui.weather2
   1: "armeabi-v7a",
   // com.kakaogames.twodin
   2: "arm64-v8a",
   // TODO: No known apps that require this
   3: "armeabi",
   // com.scottgames.fivenightsatfreddys
   4: "x86_64",
   // com.google.android.gms
   5: "mips",
}

func (p Platform) String() string {
   return Platforms[int(p)]
}

func (p *Platform) Set(s string) error {
   v, err := strconv.Atoi(s)
   if err != nil {
      return err
   }
   *p = Platform(v)
   return nil
}

type Device struct {
   // developer.android.com/guide/topics/manifest/supports-gl-texture-element
   Texture []string
   // developer.android.com/guide/topics/manifest/uses-library-element
   Library []string
   // developer.android.com/guide/topics/manifest/uses-feature-element
   Feature []string
   // developer.android.com/ndk/guides/abis
   Platform string
}
