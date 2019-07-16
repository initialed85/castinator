# castinator
Repeat a UDP unicast/multicast/broadcast to another UDP unicast/multicast/broadcast address

## Example

IPv4 Multicast on `en0` to IPv6 link-local anycast on `en5`

    ./castinator en0 239.255.137.1:6291 en5 [ff02::1%en0]:6291
