#### Synopsis

GoDSP provides API to call dsp without login and socket maintain

#### API Example
  ```
  curl -X GET "localhost:9090/v1/dsp/select?oprid=getAll&msisdn=628118003585" -H "accept: application/json"
  ```
  ##### Configuration App.conf
  ```
#dsp
dspIP     = "xxxxxx"
dspPort   = "xxxxxx"
dspUser   = "xxxxxx"
dspPwd    = "xxxxxx"
localIP   = "xxxxxx"
localPort = "xxxxxx"
dspPool = ""

oprid.getAll="<GET_SUBDATA><MSISDN>[msisdn]</MSISDN><GETCOLUMN>IMEI&amp;IMSI&amp;LAC&amp;CI&amp;PROPINSI&amp;REGIONAL&amp;BRANCH&amp;KABUPATEN&amp;PRODUCT_ID&amp;AREA&amp;SPOSNAME&amp;OSVENDOR&amp;TVENDOR&amp;TTYPE&amp;OSVERSION&amp;POINAME&amp;CLUSTER&amp;LATITUDE&amp;LONGITUDE&amp;KECAMATAN&amp;KECAMATAN&amp;KELURAHAN&amp;POINAME&amp;POILONGITUDE&amp;POILATITUDE</GETCOLUMN></GET_SUBDATA>"
#oprid.getAll="<GET_SUBDATA><MSISDN>[msisdn]</MSISDN><GETCOLUMN>IMEI</GETCOLUMN></GET_SUBDATA>"

```
#### Installation or Development

1. install golang
2. clone this project
3. go get github.com/astaxie/beego 
4. update config in 'conf/App.conf'
5. run bee run -downdoc=true -gendoc=true 
    or bee run
    or go run godb.go
6. test API http://localhost:8080/swagger/
7. dashboard & monitoring http://localhost:8088

#### API Reference
- https://beego.me


#### Tests

Describe and show how to run the tests with code examples.

#### Contributors

Let people know how they can dive into the project, include important links to things like issue trackers, irc, twitter accounts if applicable.

#### License

A short snippet describing the license (MIT, Apache, etc.)
