---
layout: post
title: Augeas the Missing Manual
description: "A guide for some JSON processing using the Augeas configuration tool"
modified: 2015-06-01
tags: [augeas]
image:
  feature: augeas.jpg
---
Lately I've been using Augeas to configure some json files, and phew! it's not been easy! So here's a little guide to help other unfortunate souls lost in the Augeas wilderness.

1) First let's start with a pretty basic json file located at <strong>/tmp/test.json</strong>:

{% highlight javascript %}
{ "a": 1 }
{% endhighlight %}

2) <strong>Load</strong>

   Loading the file with the json lens is easy enough once you know how
{% highlight plain %}
augtool> set /augeas/load/Json/lens Json.lns
augtool> set /augeas/load/Json/incl /tmp/test.json
augtool> load
augtool> print /files/tmp/test.json

/files/tmp/test.json
/files/tmp/test.json/dict
/files/tmp/test.json/dict/entry = "a"
/files/tmp/test.json/dict/entry/number = "1"
{% endhighlight %}
3) <strong>Set</strong>

   As is changing the existing property
{% highlight plain %}
augtool> set /files/tmp/test.json/dict/entry[. = "a"]/number 2
augtool> save
Saved 1 file(s)
{% endhighlight %}
  resulting in
{% highlight javascript %}
{ "a" : 2 }
{% endhighlight %}
4) <strong>Set Type</strong>

   Of course you can change this to a string or boolean property as required
{% highlight plain %}
augtool> rm /files/tmp/test.json/dict/entry[. = "a"]/number
augtool> set /files/tmp/test.json/dict/entry[. = "a"]/string hello
augtool> save
Saved 1 files(s)
{% endhighlight %}
   to create
{% highlight javascript %}
{ "a": "hello" }
{% endhighlight %}
5) <strong>Add Sub-object</strong>

   And finally adding some subobjects with properties
{% highlight plain %}
set /files/tmp/test.json/dict/entry[. = "b"] b
set /files/tmp/test.json/dict/entry[. = "b"]/dict/entry[. = "c"] c
set /files/tmp/test.json/dict/entry[. = "b"]/dict/entry[. = "c"]/dict/entry[. = "d"] d
set /files/tmp/test.json/dict/entry[. = "b"]/dict/entry[. = "c"]/dict/entry[. = "d"]/string "world"
augtool> save
Saved 1 files(s)
{% endhighlight %}
   to create an json structure with depth:
{% highlight plain %}
{
    "a": "hello",
    "b": {
        "c": {
            "d": "world"
        }
    }
}
{% endhighlight %}
