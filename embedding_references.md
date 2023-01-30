## Embedding References

References are part of a Layer 3 protocol. The most natural location for them is in a Layer 3 protocol header such as IPv4 options or IPv6 extension headers. Unfortunately, many Internet Service Providers drop options and extension headers that they deem not worthy processing. Even if they don't, network devices tend to put packets containing them on a slow processing path resulting in poor performance. The most reliable way of transferring packets with references is via a UDP tunnel. Other tunneling techniques might also be utilized.

### IPv4 Option

With IPv4 packets, references may be embedded in an IP header option.
```
    ┌───────────────┬───────────────┬───────────────┬───────────────┐
    │       0       │       1       │       2       │       3       │
    ├───────────────┼───────────────┼───────────────┼───────────────┤
    │0 1 2 3 4 5 6 7│0 1 2 3 4 5 6 7│0 1 2 3 4 5 6 7│0 1 2 3 4 5 6 7│
    ├─┬───┬─────────┼───────────────┼───────────────┼───────────────┤
    │c│cls│ opt num │    length     │           reserved            │
    ├─┴───┴─────────┼───────────────┼───────────────┼───────────────┤
    │                                                               │
    │                             sref                              │
    │                                                               │
    │                                                               │
    ├───────────────┴───────────────┴───────────────┴───────────────┤
    │                                                               │
    │                             dref                              │
    │                                                               │
    │                                                               │
    └───────────────┴───────────────┴───────────────┴───────────────┘

  Fields:

    c           - copy to fragments, set to 0

    cpt cls     - option class, set to 0

    opt num     - option number, use experimental 30 (intended new value 26)

    length      - option length in octets, set to 36

    reserved    - set to 0

    sref/dref   - source and destination references, 16 octets each
```

### IPv6 Extension Header

With IPv6 packets, references may be embedded in an IPv6 extension header.
```
    ┌───────────────┬───────────────┬───────────────┬───────────────┐
    │       0       │       1       │       2       │       3       │
    ├───────────────┼───────────────┼───────────────┼───────────────┤
    │0 1 2 3 4 5 6 7│0 1 2 3 4 5 6 7│0 1 2 3 4 5 6 7│0 1 2 3 4 5 6 7│
    ├───────────────┼───────────────┼───────────────┼───────────────┤
    │  next header  │    length     │           reserved            │
    ├───────────────┼───────────────┼───────────────┼───────────────┤
    │                            padding                            │
    ├───────────────┼───────────────┼───────────────┼───────────────┤
    │                                                               │
    │                             sref                              │
    │                                                               │
    │                                                               │
    ├───────────────┴───────────────┴───────────────┴───────────────┤
    │                                                               │
    │                             dref                              │
    │                                                               │
    │                                                               │
    └───────────────┴───────────────┴───────────────┴───────────────┘

Extension header type is experimental 254 (intended new value 145)

  Fields:

    next header - type of the next header

    length      - option length over first 8 in 8 octet increments, set to 4

    reserved    - set to 0

    padding     - set to 0

    sref/dref   - source and destination references, 16 octets each
```

### UDP Tunnel

The most reliable way of embedding references is to place them in a UDP packet. This works for both IPv4 and IPv6. In both cases, the reference is embedded in the form of a respective option or extension header. In addition, a tunnel encapsulation record is added. The order of items is as follows:

- UDP header
- Tunnel encapsulation record
- IPREF option
- Packet payload

```
    ┌───────────────┬───────────────┬───────────────┬───────────────┐
    │       0       │       1       │       2       │       3       │
    ├───────────────┼───────────────┼───────────────┼───────────────┤
    │0 1 2 3 4 5 6 7│0 1 2 3 4 5 6 7│0 1 2 3 4 5 6 7│0 1 2 3 4 5 6 7│
    └───────────────┴───────────────┴───────────────┴───────────────┘
  UDP header
    ┌───────────────┬───────────────┬───────────────┬───────────────┐
    │            src port           │            dst port           │
    ├───────────────┼───────────────┼───────────────┼───────────────┤
    │             length            │            checksum           │
    └───────────────┴───────────────┴───────────────┴───────────────┘
  Tunnel encapsulation
    ┌───────────────┬───────────────┬───────────────┬───────────────┐
    │     ttl       │  protocol     │      hops     │   reserved    │
    └───────────────┴───────────────┴───────────────┴───────────────┘
  IPREF option
    ┌───────────────┬───────────────┬───────────────┬───────────────┐
    │    option     │    length     │           reserved            │
    ├───────────────┼───────────────┼───────────────┼───────────────┤
    │                                                               │
    │                             sref                              │
    │                                                               │
    │                                                               │
    ├───────────────┼───────────────┼───────────────┼───────────────┤
    │                                                               │
    │                             dref                              │
    │                                                               │
    │                                                               │
    └───────────────┴───────────────┴───────────────┴───────────────┘


  UDP fields:

    src port    - 1045  proposed

    dst port    - 1045  same as src port

  ENCAP fields:

    ttl         - ttl copied from the original IP header

    protocol    - protcol copied from the original IP header

    hops        - number of hops detected for incoming packets

  IPREF option fields:

    option      - copy of the option or extension header first octet (ignored)

    length      - length per respective option or extension header (36 or 4)

    reserved    - set to 0

    sref/dref   - source and destination references, 16 octets each
```
