
---
title: "Gate"
linkTitle: "Gate"
date: 2024-01-04
description: >
  
---

{{< figure src="/smart-home/img/gate-network.svg" >}}

Gate Server in the Smart Home system is a key component that provides secure and convenient remote access to the smart home, even in the absence of a static IP address. This server enables users to control and monitor their home from anywhere in the world via the internet, ensuring data security and comfortable interaction.

#### Advantages of the Gate Server in the Smart Home System

1. **Security:**
    - **Data Encryption:** All data transmitted between devices and the Gate Server is encrypted to prevent unauthorized access.
    - **Authentication:** Built-in authentication mechanisms ensure that only authorized users have access to the system.

2. **Convenience of Remote Access:**
    - **No White IP:** Users can connect to the system even without a static IP address, making remote access more convenient and accessible.

3. **Configuration Flexibility:**
    - **Connect to Test Gateway:** Smart Home provides a built-in Gate client that can connect to a test gateway at https://gate.e154.ru:8443 for testing and configuring remote access.

4. **Scalability:**
    - **Support for Multiple Devices:** The Gate Server is designed to handle simultaneous requests from multiple devices, ensuring system scalability.

#### Working with the Gate Client in the Smart Home System

1. **Setting Up the Gate Client:**
    - Enable the Gate Server mode in the Smart Home system, specifying the address of the test gateway or setting up your Gate Server.
    - Configure security parameters such as encryption and authentication.

2. **Connecting to the Gate Server:**
    - The Gate client automatically establishes a secure connection with the Gate Server, without requiring a white IP address.
    - Users can use mobile applications or the web interface for remote access.

3. **Secure Remote Management:**
    - Users can control devices, monitor the home's status, and receive notifications even when away from home.

### Server Configuration

```bash
cat config.gate.json

{
  "api_http_port": 8080,
  "api_https_port": 8443,
  "api_debug": false,
  "api_gzip": true,
  "pprof": false,
  "domain": "localhost",
  "https": false,
  "proxy_timeout": 5,
  "proxy_idle_timeout": 10,
  "proxy_secret_key": ""
}
```

Properties for the `config.gate.json` configuration file:

1. **`api_http_port` (int):**
   - **Description:** Port for the HTTP API server.
   - **Example Value:** `8080`.

2. **`api_https_port` (int):**
   - **Description:** Port for the HTTPS API server.
   - **Example Value:** `8443`.

3. **`api_debug` (bool):**
   - **Description:** Enable debug mode for the API server.
   - **Example Value:** `true` (enabled).

4. **`api_gzip` (bool):**
   - **Description:** Enable Gzip compression for API requests.
   - **Example Value:** `true` (enabled).

5. **`domain` (string):**
   - **Description:** Domain name for the Gate server.
   - **Example Value:** `example.com`.

6. **`pprof` (bool):**
   - **Description:** Enable server profiling mode.
   - **Example Value:** `true` (enabled).

7. **`https` (bool):**
   - **Description:** Enable the use of Let's Encrypt for automatic SSL certificate acquisition for the specified domain.
   - **Example Value:** `true` (enabled).

8. **`proxy_timeout` (int):**
   - **Description:** Timeout for proxy connections in seconds.
   - **Example Value:** `5`.

9. **`proxy_idle_timeout` (int):**
   - **Description:** Timeout for proxy connections in seconds when there is no activity.
   - **Example Value:** `10`.

10. **`proxy_secret_key` (string):**
   - **Description:** Secret key to ensure the security of proxy connections.
   - **Example Value:** `mySecretKey`.

These parameters provide flexible control over the settings of the Gate server, including security, operating modes, and the use of SSL certificates via Let's Encrypt.

### Server Launch

The Gate server is integrated into the smart-home system as a separate mode by specifying the `gate` argument during startup.

```bash
./smart-home help gate
Organization of remote access without white IP

Usage:
  server gate [flags]

Flags:
  -h, --help   help for gate


./smart-home gate

  ___      _
 / __|__ _| |_ ___
| (_ / _' |  _/ -_)
 \___\__,_|\__\___|


INFO	gate            	server/gate_server.go:93 >	Started ...
INFO	gate            	server/server.go:117 >	server started at :8080
```
