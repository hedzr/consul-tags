# consul

`~/Library/LaunchAgents/homebrew.mxcl.consul.plist`:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
        <key>KeepAlive</key>
        <true/>
        <key>Label</key>
        <string>homebrew.mxcl.consul</string>
        <key>LimitLoadToSessionType</key>
        <array>
                <string>Aqua</string>
                <string>Background</string>
                <string>LoginWindow</string>
                <string>StandardIO</string>
                <string>System</string>
        </array>
        <key>ProgramArguments</key>
        <array>
                <string>/opt/homebrew/opt/consul/bin/consul</string>
                <string>agent</string>
                <string>-dev</string>
                <string>-bind</string>
                <string>127.0.0.1</string>
        </array>
        <key>RunAtLoad</key>
        <true/>
        <key>StandardErrorPath</key>
        <string>/opt/homebrew/var/log/consul.log</string>
        <key>StandardOutPath</key>
        <string>/opt/homebrew/var/log/consul.log</string>
        <key>WorkingDirectory</key>
        <string>/opt/homebrew/var</string>
</dict>
</plist>
```

`/opt/homebrew/var/log/consul.log`:

```bash
==> Starting Consul agent...
               Version: '1.16.2'
            Build Date: '1970-01-01 00:00:01 +0000 UTC'
               Node ID: '11d3461b-af32-f2d2-c822-cae1537090fc'
             Node name: 'TODD-14.local'
            Datacenter: 'dc1' (Segment: '<all>')
                Server: true (Bootstrap: false)
           Client Addr: [127.0.0.1] (HTTP: 8500, HTTPS: -1, gRPC: 8502, gRPC-TLS: 8503, DNS: 8600)
          Cluster Addr: 127.0.0.1 (LAN: 8301, WAN: 8302)
     Gossip Encryption: false
      Auto-Encrypt-TLS: false
           ACL Enabled: false
     Reporting Enabled: false
    ACL Default Policy: allow
             HTTPS TLS: Verify Incoming: false, Verify Outgoing: false, Min Version: TLSv1_2
              gRPC TLS: Verify Incoming: false, Min Version: TLSv1_2
      Internal RPC TLS: Verify Incoming: false, Verify Outgoing: false (Verify Hostname: false), Min Version: T
LSv1_2
```

## consul.json

```json
{
    "ui_dir": "/Users/hz/Downloads/data/consul.data/www"
}
```
