# pad

[Message Padding](https://en.wikipedia.org/wiki/Padding_(cryptography))

A small utility allowing users to pad their messages to a fixed size, which they may always use, to make it harder for Eve to guess the content length of a given encrypted message.


Usage: **pad -p 4096 < infile > outfile**, to unpad: **pad -u < infile > outfile**.

```
+-----------------------------------------------------------------------+
|<-------------------------- padsize bytes ---------------------------->|
+-----------------------------------------------------------------------+
|<--- content bytes --->|<--- padding bytes --->|<--- uint64t value --->|
+-----------------------------------------------------------------------+

uint64t value is a little endian binary number representing byte_length(content).
```

