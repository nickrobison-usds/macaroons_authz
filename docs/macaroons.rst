Macaroons Format
================

The Macaroons format is kind of confusing, so we'll try to explain it a bit better.

Library support is really varied, here's what we know so far:

========== ============= ============ ========= =========== =========
Language   Library       V1  (Binary) V1 (JSON) V2 (Binary) V2 (JSON)
========== ============= ============ ========= =========== =========
Go         go-macaroons  Yes          Yes       Yes         Yes
Java       jmacaroons    Yes          No        No          No
Javascript js-macaroons  No           No        Yes         Yes
C/C++      libmacaroons  Yes          No        Yes         Yes
========== ============= ============ ========= =========== =========


Caveats
-------

There are first party and third-party caveats.
First party are easy, third-party.... not so much.


Third-party caveats
~~~~~~~~~~~~~~~~~~~

Third party caveats have two different versions.

Note: Most macaroons libraries do not support automatically encrypting the caveat data! You'll need to do it yourself using `NaCL`_.

A third-party caveat is essentially a first-part caveat that is encrypted using a *NaCL box* and sent to the third party, along with the senders public key, which is used for decryption.

**V2**

The V2 caveats have the following format:

```
version 2 or 3 [1 byte]
first 4 bytes of third-party Curve25519 public key [4 bytes]
first-party Curve25519 public key [32 bytes]
nonce [24 bytes]
encrypted secret part [rest of message]
```

Note: The V3 caveats are only supported by the `macaroon-bakery`_ library for Go.

What's a bit tricky is that you send along the first 4 bytes of the public key used to encrypt the data, along with the entirety of *your* public key, which is used by the third party to decrypt the data.
The actual data that is encrypted is laid out as follows:

```
version 2 or 3 [1 byte]
root key length [n: uvarint]
root key [n bytes]
predicate [rest of message]
```

The root key is then used to decrypt the predicate using a *NaCL secretbox*.

The entire message is then base64 encoded and sent along to the third party.


.. _NaCL: https://nacl.cr.yp.to
.. _macaroon-bakery: https://github.com/go-macaroon-bakery/macaroon-bakery
