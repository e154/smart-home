---
title: "http"
linkTitle: "http"
date: 2021-10-20
description: >

---

In the **Smart Home** project, there is a capability to perform arbitrary synchronous HTTP requests to external
resources.

The `http` object allows you to make synchronous HTTP requests to external resources, such as API services, and receive
responses. You can use this method to integrate with other systems and retrieve or send data through the HTTP protocol
in your **Smart Home** project.

Supported methods include:

* GET
* POST
* PUT
* DELETE

{{< alert color="success" >}}This function is available in any system script.{{< /alert >}}

----------------

For this purpose, the corresponding methods are available:

### GET Request

```coffeescript
response = http.get(url)
```

| Parameter | Description                 |
|-----------|-----------------------------|
| url       | The request URL             |
| response  | The response of the request |

----------------

### POST Request

```coffeescript
response = http.post(url, body)
```

| Parameter | Description                 |
|-----------|-----------------------------|
| url       | The request URL             |
| body      | The request body            |
| response  | The response of the request |

----------------

### Headers Request

```coffeescript
response = http.headers(headers).post(url, body)
```

| Parameter | Description                 |
|-----------|-----------------------------|
| headers   | The request headers         |
| url       | The request URL             |
| body      | The request body            |
| response  | The response of the request |

----------------

### пример кода

```coffeescript
# auth
# ##################################

res = http.digestAuth('user', 'password').download(uri);

res = http.basicAuth('user', 'password').download(uri);

res = http.download(uri);


# GET http
# ##################################

res = http.get("%s")
if res.error
  return
p = JSON.parse(res.body)


# POST http
# ##################################

res = http.post("%s", {'foo': 'bar'})
if res.error
  return
p = JSON.parse(res.body)


# PUT http
# ##################################

res = http.put("%s", {'foo': 'bar'})
if res.error
  return
p = JSON.parse(res.body)


# GET http + custom headers
# ##################################

res = http.headers([{'apikey': 'some text'}]).get("%s")
if res.error
return
p = JSON.parse(res.body)

# DELETE http
# ##################################

res = http.delete("%s")
if res.error
  return
p = JSON.parse(res.body)

```
