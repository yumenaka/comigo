<h2 align="center">
  <img src="https://raw.githubusercontent.com/yumenaka/comi/master/icon.ico" alt="ComiGo: simple comic reader" width="200">
  <br>ComiGo: simple comic reader<br>
</h2>

![Windows example](https://www.yumenaka.net/wp-content/uploads/2020/08/sample.gif "Windows example")  

## Supported platforms：

- Windows(386\64)
- macOS
- Linux（amd64、armv7\armv8）
- and More(You can compile it yourself).

### Install：
Download [the latest release for your system](https://github.com/yumenaka/comi/releases ) 

### Usage：
```
comi [flags] file_or_dir

Set the web service port (default is 1234): 
comi -p 2345 book.zip 

Without opening a browser (windows):
comi -o=false book.zip
   
For local browsing, not available for LAN:
comi -l book.zip 

webp transfer, requiring webp-server to work with: 
comi -w book.zip 

Multiple parameters:
comi -w -q 70 --frpc  --token aX4457d3O -p 23455 --frps-addr sh.example.com test.zip

Flags:
      --check-image                Analyze the image resolution on the server side (default true)
  -c, --config string              Specify configuration file
      --debug                      Debug Mode
  -d, --disable-lan                Disable LAN sharing
  -f, --frpc                       Enable frp reverse proxy
      --frpc-command string        frpc command, or frpc executable file path, the default isfrpc (default "frpc")
      --frps-addr string           frps-addr, frpc need (default "frps.example.com")
      --frps-port int              frps server_port, frpc need (default 7000)
      --frps-random-remote         frp remote random port(40000~50000) (default true)
      --frps-remote-port int       frp remote port, if set to -1, the port is the same as the local (default 50000)
  -h, --help                       help for comi
      --host string                Custom domain name
      --log                        Enable logFile
  -m, --max-depth int              Maximum search depth (default 1)
      --min-image-num int          There are at least a few media file before it is considered a comic compression package (default 1)
  -o, --open-browser               Open the browser at the same time，windows=true
  -p, --port int                   Service port (default 1234)
      --print-all-ip               Print all available network card ip
      --sketch_count_seconds int   Countdown seconds in sketch mode (default 90)
      --sort string                Image rearrangement rules (none,name,time) (default "none")
  -t, --template string            Default page template(scroll,single,double,sketch) (default "scroll")
      --token string               token, frpc need (default "token_secretSAMPLE")
  -v, --version                    version for comi
  -w, --webp                       webp transmission, webp-server is required
      --webp-command string        webp-server command, or webp-server executable file path (default "webp-server")
  -q, --webp-quality int           webp compression quality (default 85)
```

### Special Thanks：
[mholt](https://github.com/mholt)  、[spf13](https://github.com/spf13)  [disintegration](https://github.com/disintegration)   、 [Baozisoftware ](https://github.com/Baozisoftware) 、 [markbates](github.com/markbates/pkger)  and more。

## License

This software is released under the GPL-3.0 license.
