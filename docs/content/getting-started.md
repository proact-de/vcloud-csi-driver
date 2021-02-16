---
title: "Getting Started"
date: 2021-02-16T00:00:00+00:00
anchor: "getting-started"
weight: 10
---

## Installation

TBD

## Configuration

`VCLOUD_CSI_LOG_LEVEL`
: Set logging level, defaults to `info`

`VCLOUD_CSI_LOG_PRETTY`
: Enable pretty logging, defaults to `true`

`VCLOUD_CSI_LOG_COLOR`
: Enable colored logging, defaults to `true`

`VCLOUD_CSI_NODENAME`
: Name of the node running on, this is a **required** parameter

`VCLOUD_CSI_HREF`
: URL to access vCloud Director API, this is a **required** parameter

`VCLOUD_CSI_INSECURE`
: Skip SSL verify for vCloud Director, defaults to `false`

`VCLOUD_CSI_USERNAME`
: Username for vCloud Director, this is a **required** parameter

`VCLOUD_CSI_PASSWORD`
: Password for vCloud Director, this is a **required** parameter

`VCLOUD_CSI_ORG`
: Organization for vCloud Director, this is a **required** parameter

`VCLOUD_CSI_VDC`
: VDCs for vCloud Director, this is a **required** parameter

`VCLOUD_CSI_ENDOINT`
: Path to unix socket endpoint, defaults to `unix:///csi/csi.sock`
