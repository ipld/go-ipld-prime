gengo
=====

This package contains a codegenerator for emitting Golang source code
for datastructures based on IPLD Schemas.

It is at present a partially complete proof-of-concept.  Use at your own risk.

There is not yet a user-facing CLI; you have to write code to use it.

Check out the [HACKME](HACKME.md) document for more info about the internals,
how they're organized, and how to hack on this package.


completeness
------------

Legend:

- `✔` - supported!
- `✘` - not currently supported.
- `⚠` - not currently supported -- and might not be obvious; be careful.
- `-` - is not applicable
- `?` - feature definition needed!  (applies to many of the "native extras" rows -- often there's partial features, but also room for more.)
- ` ` - table is not finished, please refer to the code and help fix the table :)

| feature                        | accessors | builders |
|:-------------------------------|:---------:|:--------:|
| structs                        |    ...    |    ...   |
| ... type level                 |     ✔     |     ✔    |
| ... native extras              |     ?     |     ?    |
| ... map representation         |     ✔     |     ✔    |
| ... ... including optional     |     ✔     |     ✔    |
| ... ... including renames      |     ✔     |     ✔    |
| ... ... including implicits    |     ⚠     |     ⚠    |
| ... tuple representation       |     ✘     |     ✘    |
| ... ... including optional     |           |          |
| ... ... including renames      |           |          |
| ... ... including implicits    |           |          |
| ... stringjoin representation  |     ✔     |     ✔    |
| ... ... including optional     |     -     |     -    |
| ... ... including renames      |     -     |     -    |
| ... ... including implicits    |     -     |     -    |
| ... stringpairs representation |     ✘     |     ✘    |
| ... ... including optional     |           |          |
| ... ... including renames      |           |          |
| ... ... including implicits    |           |          |
| ... listpairs representation   |     ✘     |     ✘    |
| ... ... including optional     |           |          |
| ... ... including renames      |           |          |
| ... ... including implicits    |           |          |

| feature                        | accessors | builders |
|:-------------------------------|:---------:|:--------:|
| lists                          |    ...    |    ...   |
| ... type level                 |     ✔     |     ✔    |
| ... native extras              |     ?     |     ?    |
| ... list representation        |     ✔     |     ✔    |

| feature                        | accessors | builders |
|:-------------------------------|:---------:|:--------:|
| maps                           |    ...    |    ...   |
| ... type level                 |     ✔     |     ✔    |
| ... native extras              |     ?     |     ?    |
| ... map representation         |     ✔     |     ✔    |
| ... stringpairs representation |     ✘     |     ✘    |
| ... listpairs representation   |     ✘     |     ✘    |

| feature                        | accessors | builders |
|:-------------------------------|:---------:|:--------:|
| unions                         |    ...    |    ...   |
| ... type level                 |     ✘     |     ✘    |
| ... keyed representation       |     ✘     |     ✘    |
| ... envelope representation    |     ✘     |     ✘    |
| ... kinded representation      |     ✘     |     ✘    |
| ... inline representation      |     ✘     |     ✘    |
| ... byteprefix representation  |     ✘     |     ✘    |

| feature                        | accessors | builders |
|:-------------------------------|:---------:|:--------:|
| strings                        |     ✔     |     ✔    |
| bytes                          |     ✔     |     ✔    |
| ints                           |     ✔     |     ✔    |
| floats                         |     ✔     |     ✔    |
| bools                          |     ✔     |     ✔    |
| links                          |     ✔     |     ✔    |

| feature                        | accessors | builders |
|:-------------------------------|:---------:|:--------:|
| enums                          |    ...    |    ...   |
| ... type level                 |     ✘     |     ✘    |
| ... string representation      |     ✘     |     ✘    |
| ... int representation         |     ✘     |     ✘    |
