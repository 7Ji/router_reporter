Simple go program to be used on the router (mainly OpenWrt) to listen on a specific port and report the IP address on a specific interface to multiple clients in a set interval.  
简单的用在路由器（主要是OpenWrt）上的用于在特定端口上监听并以恒定间隔向多个客户端广播的程序

I wrote this program so I could have a QQ bot to run on a internal server to watch on the router's public IP and notify me if its IP has changed, without relying on external IP reporting services (e.g. ipw.cn).   
写这个程序的目的是为了在内网服务器上跑QQ机器人，监视路由器的公网IP，并且在IP变动的情况下来汇报，而不依赖外部的IP汇报服务（比如ipw.cn）

The file `openwrt_router_reporter.procd.service` should be modified with `iface` and `listen` set according to your actual network configuration, and placed as `/etc/init.d/router_reporter` and with mod `755` to be enabled on an OpenWrt device
文件`openwrt_router_reporter.procd.service`里面的`iface`和`listen`根据具体网络配置修改后，应当被放置为`/etc/init.d/router_reporter`并且设置权限为`755`来在OpenWrt设备上开机自启

The command to compile the program to be run on a MIPSLE softfloat platform (e.g. MT7621, the SoC used by Phicomm K2P and Xiaomi Redmi AC2100):  
编译在一个小端MIPS软浮点平台（比如说斐讯K2P和红米AC2100上的SoC MT7621）上运行的程序的命令
```
GOARCH=mipsle GOMIPS=softfloat go build -ldflags "-s -w -buildid=" -trimpath
```
*Go can cross-compile easily without extra setting about the toolchain, the above command can be simply run on an x86-64 host natively
Go可以简单地不经额外工具链设置就能交叉编译，上面的命令可以在x86-64的机子上原生运行*
```
$ uname -a
Linux laptop7ji 6.0.8-zen1-1-zen #1 ZEN SMP PREEMPT_DYNAMIC Thu, 10 Nov 2022 21:14:22 +0000 x86_64 GNU/Linux
$ GOARCH=mipsle GOMIPS=softfloat go build -ldflags "-s -w -buildid=" -trimpath
$ file router_reporter
router_reporter: ELF 32-bit LSB executable, MIPS, MIPS32 version 1 (SYSV), statically linked, stripped
```

To test if it is running correctly, simple use `nc`:  
要检查是否正常运行，可以简单地使用`nc`：
```
nc [listen host] [listen port]
```
Expected output:  
正常的输出：
```
$ nc 127.0.0.1 7777
192.168.7.25/24
192.168.7.25/24
192.168.7.25/24
192.168.7.25/24
```
