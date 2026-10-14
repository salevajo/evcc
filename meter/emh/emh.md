# EMH metering Smart Gateway

## Local access using LAN

Connect a free local Ethernet port (e.g. `eth1`) of your evcc host to the Ethernet port on your EMH device.

If required, install `nmcli` using the `network-manager` package.

```shell
sudo apt install network-manager
```

Configure a local connection using the default CIDR value `192.168.33.10/24` and set connection up.

```shell
sudo nmcli connection modify "eth1" ipv4.method manual ipv4.addresses 192.168.33.10/24   ipv4.never-default yes ipv6.never-default yes
sudo nmcli connection up "eth1"
```

Add `smgw.local` to your local `/etc/hosts` using the configured IP address of your EMH device.

```/etc/hosts
192.168.33.10 smgw.local
```

## Username and password to access EMH

The individual username and password will be provided by your local grid provider, there are no default values to access any data.

The username, as an example _EWE Netz_ is using a long numeric value, and password are required for HTTP basic auth.

Some URLs need the smart meter ID, like `https://smgw.local/json/metering/origin/1efr0123456789`. The ID `1efr0123456789` is an individual value for each smart meter.

## Data provided

### Current values

While using the web based interface, you can get HTML pages with informations about the current smart meter values like grid total [kWh], voltage for every phase [V], utilization [W], power [A].

A request to [](https://smgw.local/json/metering/origin/1efr0123456789/extended) provides a JSON representation for current values like in the table below. JSON file see [here](start.exteded.json). In this file, the field `logical_name` is a dot delimited list of three parts, _register ID_, _smart meter ID_ and the suffix `sm`.

| Timestamp           |	Register     | Value with unit
| ------------------- | ------------ | ----------------
| 2026-10-13 17:31:20	| 0100010800ff | 1.918,637000 kWh
| 2026-10-13 17:31:20	| 0100020800ff | 0,000000 kWh
| 2026-10-13 17:31:20	| 01000e0700ff | 50,000000 Hz
| 2026-10-13 17:31:20	| 0100100700ff | 1.433,000000 W
| 2026-10-13 17:31:20	| 01001f0700ff | 1,700000 A
| 2026-10-13 17:31:20	| 0100200700ff | 232,100000 V
| 2026-10-13 17:31:20	| 0100240700ff | 334,000000 W
| 2026-10-13 17:31:20	| 0100330700ff | 1,530000 A
| 2026-10-13 17:31:20	| 0100340700ff | 231,700000 V
| 2026-10-13 17:31:20	| 0100380700ff | 295,000000 W
| 2026-10-13 17:31:20	| 0100470700ff | 3,510000 A
| 2026-10-13 17:31:20	| 0100480700ff | 230,300000 V
| 2026-10-13 17:31:20	| 01004c0700ff | 802,000000 W
| 2026-10-13 17:31:20	| 0100510701ff | 240,000000 °
| 2026-10-13 17:31:20	| 0100510702ff | 120,000000 °
| 2026-10-13 17:31:20	| 0100510704ff | 340,000000 °
| 2026-10-13 17:31:20	| 010051070fff | 340,000000 °
| 2026-10-13 17:31:20	| 010051071aff | 353,000000 °

### Event log

At [](https://smgw.local/conlog) a event log can be fetched on a monthly base having [0](https://smgw.local/conlog/0) for the current month and counting up like [1](https://smgw.local/conlog/1) and so forth for all available month. A JSON file for each event log is available using [](https://smgw.local/json/conlogbook/0/extended) for the current month. See [here](events.extended.json) for an example.

### Meterlist

On page [](https://smgw.local/meterlist) you can find individual links to active and archived meter tarif data in different granularity, e.g. monthly or data sent every 15 minutes. For each link you can fetch detailed data using an individual logical name [](https://smgw.local/tafrecords/{logical_name}) and for every available timeframe in each section a [meter_records.csv](https://smgw.local/tafrecords/{logical_name}/1/download) or a JSON formatted file [meter_records.json](https://smgw.local/json/metering/derived/{logical_name}/billingperiods/1/extended).


## Validation certificate

You can get a [certificate.pem](https://smgw.local/download/certificate.pem) to check measurements provided by the smart meter.

## Perform self-test

A selftest initiated using [](https://smgw.local/selftest) will be available after a short period of time in the event logs.

# evcc Integration

## Configuration

| variable     | default value
| ------------ | ------------------------------------------------------------------
| interface    | `eth1`
| ipv4         | `192.168.33.10/24`
| host         | `smgw.local`
| id           | `1efr0123456789`
| logical_name | `01234e012345.ewesup1000000123456taf01000012345678000006.sm`
| poll_seconds | `300`

## Provided to evcc

- Stores polled values in a local table to store historical data about your power consumption.
- Import meterlist data
