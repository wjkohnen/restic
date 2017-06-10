Local Cache
===========

In order to speed up certain operations, restic manages a local cache of data.
This document describes the data structures for the local cache. This document
describes the cache with version 1.

Versions
--------

The cache directory is selected according to the `XDG base dir specification
<http://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html>`__. It
contains a file named `version`, which contains a single ASCII integer line
that stands for the current version of the cache. If a lower version number is
found the cache is recreated with the current version. If a higher version
number is found the cache is ignored and left as is.

Snapshots and Indexes
---------------------

Snapshot and index files are cached in the sub-directories ``snapshots`` and
``index``, as read from the repository.

Keys
----

Keys are never cached locally. The sub-directory ``key`` contains an empty file
which has the same name as the key file that was successfully used the last
time to open the repository. On a subsequent run, restic tries this key file
first (if it still exists).

Blobs
-----

The sub-directory ``Blobs`` stores Blobs combined to larger files. The directory
contains one level of sub-directories, named ``00`` to ``ff``. These
sub-directories store files with names between two and 64 hexadecimal
characters. A file may only store Blobs whose SHA-256 hash starts with the file
name. For example, consider the following cache directory for the repo with the
ID starting with ``263a9c7bbc``:

::

    $ tree .cache/restic
    .cache/restic
    ├── 263a9c7bbc494ab1c7aa67fe74f8794c309e2f91c0803884b0c31670bbcd4a54
    │   └── blobs
    │       ├── 05
    │       │   ├── 0501020504
    │       │   └── 0506fffefdfc
    │       └── c0
    │           └── c0c1ff
    └── version

For example, the file ``05/0501020504`` may only store Blobs whose SHA-256 hash
starts with ``0501020504``.

A Blob cache file consists of a sequence of (encrypted) Blobs. The Blobs are
encrypted with the repository master key. Each encrypted Blob is prepended by a
32 bit unsigned integer (``uint32``) in little-endian encoding. The most
significant byte has the following meaning:

+--------+----------------+
| Type   | Meaning        |
+========+================+
| 0x00   | data blob      |
+--------+----------------+
| 0x01   | tree blob      |
+--------+----------------+
| 0xff   | cache header   |
+--------+----------------+

The other three bytes are the length of the encrypted Blob that follows. This
structure allows iterating through the Blobs sequentially. In the following,
the notation ``LengthType(x)`` is used for encoding the type and the length as
a four byte unsigned integer in little-endian encoding.

On startup, the Blob cache files need to be read so that restic knows which
Blobs have been cached locally. This is sped up by loading an (encrypted) cache
header from the file. It is always located at the end of the file, the length
is repeated in the last four bytes so that restic seek to the offset -4 and
read the length of the header.

So a typical Blob file may look like this:

::

    LengthType(Blob1_Data) || Blob1_Data ||
    LengthType(Blob2_Data) || Blob2_Data ||
    [...]
    LengthType(BlobN_Data) || BlobN_Data ||
    LengthType(Header+LengthType(Header)) || Header || LengthType(Header)

This structure can be read either sequentially. In this case, the header length
field contains the length of the header itself, plus the four bytes for the
additional length field at the end. So restic needs to subtract four bytes from
the length to get the length of the encrypted cache header data. When read from
the end, no subtraction is needed.

Blob cache files can be appended, so that the files don't need to be rewritten
from scratch when new Blobs are to be cached. A new header is appended to the
file, so the header at the end of the file is the current one.

The example Blob file introduced above may then look like this:

::

    LengthType(Blob1_Data)   || Blob1_Data   ||
    LengthType(Blob2_Data)   || Blob2_Data   ||
    [...]
    LengthType(BlobN_Data)   || BlobN_Data   ||
    LengthType(Header+LengthType(Header)) || Header || LengthType(Header)
    LengthType(BlobN+1_Data) || BlobN+1_Data ||
    LengthType(BlobN+2_Data) || BlobN+2_Data ||
    [...]
    LengthType(BlobM_Data)   || BlobM_Data   ||
    LengthType(Header+LengthType(Header_New)) || Header || LengthType(Header_New)

Iterating sequentially through the file, restic finds two headers and uses only
the last one. When looking from the end of the file, the last header can be
determined easily. The last header covers the complete file.

Header Format
~~~~~~~~~~~~~

The header (after decryption) is a sequence of JSON documents each describing
the offset, length, plaintext SHA-256 hash and type of a Blob in the cache
file. A sample cache header may look like this:

.. code:: json

    {
      "id": "3ec79977ef0cf5de7b08cd12b874cd0f62bbaf7f07f3497a5b1bbcc8cb39b1ce",
      "type": "data",
      "offset": 0,
      "length": 25
    }
    {
      "id": "9ccb846e60d90d4eb915848add7aa7ea1e4bbabfc60e573db9f7bfb2789afbae",
      "type": "tree",
      "offset": 38,
      "length": 100
    }
