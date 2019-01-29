# tbln

TBLN(Table notation) format read/write library.

## Specification

1. One line ends with a line feed
2. Value consists of multiple columns.
3. Data starts with "| " (vertical bar + space).
4. Values are separated by " | "  (space + vertical bar + space).
5. If you want to include "|" in the value, it is "||".
   Increase | after | one by one.
6. Otherwise, all values are taken.

## Example

| column1 | column2 | column3 |

| column| | |column2 | co|umn3 |

| co || lumn | co ||| lumn | co |||| lumn |

| co\nlumn1 | column2\n | column3\n\n |
