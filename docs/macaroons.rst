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
========== ============= ============ ========= =========== =========


Caveats
-------

There are first party and third-party caveats.
First party are easy, third-party.... not so much.


Third-party caveats
~~~~~~~~~~~~~~~~~~~

Third party caveats have two different versions.

Note: Most macaroons libraries do not support automatically encrypting the caveat data! You'll need to do it yourself using `NaCL`_

.. _NaCL: https://nacl.cr.yp.to
