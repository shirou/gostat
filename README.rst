===============
GoStat
===============

GoStat is monitoring systems which is written in golang. GoStat allow
you to view some of your system resources in real-time.

note: This is my golang practicing code. Don't use in production.

Installation / Usage
----------------------------

::

  % mkdir gostat
  % cd gostat
  % export GOPATH=`pwd`
  % go get bitbucket.org/r_rudi/gostat

::

  % ./bin/gostat
  time:2013-05-28 20:53:50.8447712 +0900 JST      tag:aio aio:0
  time:2013-05-28 20:53:50.844889839 +0900 JST    tag:load load1:0.00      load5:0.00      load15:0.00
  time:2013-05-28 20:53:50.845009392 +0900 JST    tag:memory usage MemFree:254008  Buffers:104108  Cached:493348
  time:2013-05-28 20:53:51.045134834 +0900 JST    tag:cpu usr:0.00 sys:0.00        idl:0.00        wai:100.00      hiq:0.00 siq:0.00        stl:0.00

Options
+++++++++++++

- i: interval time. default is 0 and not loop.
- o: output format ("ltsv", "csv", "whitespace", "mqtt", "http")

If you choose mqtt or http, you need specify server url.

- http

  ::

    % ./gostat -o http http://example.com/push

- mqtt

  ::

    % ./gostat -o mqtt mqtt.example.com:1833

  or if you want to specify topic,

  ::

    % ./gostat -o mqtt mqtt.example.com:1833 r_rudi/gostat/linux

  Default topic is "gostat"

Features
-----------

modules
++++++++


- load avg
- cpu
- memory
- aio

output format
++++++++++++++++

- LTSV (http://ltsv.org) to stdout
- CSV to stdout
- WhiteSpace separeted to stdout
- HTTP

  - send http post in 'json' values. it's suit for flueentd.

- MQTT (http://mqtt.org)

Limitation
----------

only work on the Linux OS.

The giants on whose shoulders this works stands
----------------------------------------------------

- dstat: http://dag.wieers.com/home-made/dstat/

License
------------------

MIT License
