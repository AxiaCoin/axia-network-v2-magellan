// (c) 2021, AXIA Systems, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package cfg

const defaultJSON = `{
  "networkID": 1,
  "logDirectory": "/tmp/magellan/logs",
  "listenAddr": ":8080",
  "chains": {},
  "services": {
    "db": {
      "dsn": "root:password@tcp(127.0.0.1:3306)/magellan_dev",
      "driver": "mysql"
    }
  }
}`
