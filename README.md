## Wordladder

A wordladder is a type of word puzzle. Per Wikipedia:

> Word ladder (also known as Doublets, word-links, change-the-word puzzles, paragrams, laddergrams, or Word golf) is a word game invented by Lewis Carroll. A word ladder puzzle begins with two words, and to solve the puzzle one must find a chain of other words to link the two, in which two adjacent words (that is, words in successive steps) differ by one letter.

As an example, to get from HEAD to TAIL:

```
HEAD
HEAL
HEIL
HAIL
TAIL
```

## API

Here is a simple wordladder API written in Go. Not included: .msgpack graphs.

### Words

##### PATH

Return the shortest path between two words.

```
GET /path/{start}/{end}
```

Example output:
```
{
    "length": 5,
    "path": [
        "head",
        "heal",
        "heil",
        "hail",
        "tail"
    ]
}
```

##### LONGPATH

Return the longest path between two words.

```
GET /path/{start}/{end}
```

Example output:
```
{
    "length": 1962,
    "path": [
        "head",
        "yead",
        "yerd",
        "yird",
        ...
        "noel",
        "noil",
        "roil",
        "rail",
        "tail"
    ]
}
```

##### NEIGHBORS

Return the all the neighbors (words which are 1 edit away) of a given word.

```
GET /neighbors/{word}
```

Example output:
```
{
    "length": 17,
    "neighbors": [
        "head",
        "bead",
        "dead",
        "heal",
        "heap",
        "hear",
        "heat",
        "heed",
        "heid",
        "held",
        "hend",
        "herd",
        "lead",
        "mead",
        "read",
        "tead",
        "yead"
    ]
}
```

### Graphs

##### NEIGHBORS

Return the stats for a graph of a given word length size.

```
GET /stats/{length}
```

Example output:
```
{
    "nodes": 5525,
    "edges": 78707,
    "most_connected": {
        "word": "tats",
        "number_of_connections": 41
    },
    "singletons": {
        "words": [
            "hwyl",
            "euoi",
            "omov",
            "ngai",
            "expo",
            "djin",
            "waac",
            "evil",
            "jehu",
            "zyga",
            "odor",
            "occy",
            "pfft",
            "ombu",
            "yunx",
            "epha",
            "jeux",
            "ahoy",
            "envy",
            "aesc",
            "eevn",
            "elhi",
            "imam",
            "exam",
            "ossa",
            "ecru",
            "epee",
            "abri",
            "mzee",
            "mwah",
            "myxo",
            "gyny",
            "asci",
            "upta",
            "kiwi",
            "isit",
            "adaw",
            "ankh",
            "odso",
            "ovum",
            "adze",
            "khor",
            "huhu",
            "acai",
            "lwei",
            "pruh",
            "yebo",
            "enuf",
            "oppo",
            "exul",
            "ygoe",
            "zoic",
            "ogam",
            "erhu",
            "vuln",
            "ekka",
            "ughs",
            "azym",
            "dzho"
        ],
        "count": 59
    }
}
```

##### WORDS

Return all the words of a given length.

```
GET /words/{length}
```

Example output:
```
{
    "count": 5525,
    "words": [
        "tora",
        "long",
        "kati",
        "wawa",
        "awry",
        "fras",
        "wonk",
        "yuck",
        "pech",
        "trim",
        "quep",
        "gule",
        "bize",
        "alod",
        "thus",
        "dell",
        ...
        "aged",
        "gabs",
        "fail",
        "leir",
        "jibe",
        "pees",
        "redd",
        "sial",
        "rias",
        "albe",
        "sync",
        "lich"
    ]
}
```
