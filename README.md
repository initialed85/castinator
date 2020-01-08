# castinator
Repeat a UDP unicast/multicast/broadcast to another UDP unicast/multicast/broadcast address

## Example

IPv4 Multicast on `en0` to IPv6 link-local anycast on `en5`

    ./castinator \
      -leftIntfcName en0 \
      -leftUDPListenAddr 239.192.137.1:6291 \
      -leftUDPSendAddr 239.192.137.1:6291 \
      -rightIntfcName en5 \
      -rightUDPListenAddr [ff02::1%en0]:6291 \
      -rightUDPSendAddr [ff02::1%en0]:6291
