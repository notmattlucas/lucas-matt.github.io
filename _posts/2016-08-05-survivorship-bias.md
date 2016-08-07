---
layout: post
title: Survivorship Bias and Negotiating Tech Hype
description: "The role of survivorship bias and
web-scale/high-performance envy when software developers choose
new tech."
modified: 2016-03-04
tags: [thinking, statistics, development]
image:
  feature: whaam.jpg
---

## Bombers with Bullet Holes

During the dark days of World War II the American military presented
their Statistical Research Group with a problem. They wanted to add additional armour to
their bombers, but clearly couldn't put the armour everywhere because
of the additional weight it would add to the planes. The group was
tasked with working out how much armor to allocate to the various regions
of the aircraft to maximize defense whilst minimizing any effect on
fuel consumption and agility.

Engineers inspected a number of bombers that had seen some action. These
planes had bullet holes that were distributed mainly across the wings
and body. Comparatively the engines and cockpit had much less
damage. This had lead the commanders to make the obvious, but foolish, conclusion that
they should enhance armour on areas that had been hit most
frequently, namely the fuselage and wings. 

One of the many geniuses of the group, Abraham Wald, realised that they
were looking at the problem from completely the wrong angle. It's not
that planes weren't being hit as frequently on the engines and cockpit,
but rather those that had been never returned to tell the tale! These
were the parts of the airplane that needed enhancement, not the areas
that could take a battering and still survive.

## Survivorship Bias

How do you really evaluate a success when the failures are nowhere
to be seen?

{:refdef: style="text-align: center;"}
![Survive]({{ site.url }}/images/survive.jpg){:
.center-image width="85%" }
{: refdef}

Countless articles, books and documentaries have been produced about
successful people and how to capture the principles of their success to
improve your own fortune.

Consider Steve Jobs - frequently heralded as a one of the greatest
geniuses of our time - how do we emulate his success? Clearly dropping
out of college, spending time at meditation retreats and starting a
business from your parent's garage is the way to go. But what about
the hundreds of thousands of budding Apple founders for whom this
strategy never quite worked out?

Books aren't usually written about failed enterprises,
just the rare, billion-dollar, success stories.

## Choosing Tech Thoughtfully

As well as being responsible for the latest diet fads and the exaggerated
performance of mutual funds, we can see this
bias lurking in certain corners of the software development world.

{:refdef: style="text-align: center;"}
![Thinking]({{ site.url }}/images/homer_thought.jpg){:
.center-image width="75%" }
{: refdef}

How often do you see a wave of enthusiasm for the next high
throughput NoSQL
system or a push for complicated elastic scaling technology? Companies
such as Twitter and Netflix present their wild successes but we don't usually see qualifications on the size and
scale of the teams implementing these solutions.

It's worth keeping in
mind the potential for a mass of silent teams. Struggling under the
weight of overpowered, over-engineered, "web-scale" technologies
inspired by the industry front-runners. Most of us mere mortals just don't
have the resources, skills, or (most importantly) even the need for
such high class deployments. Most of the time it's just better to keep
things simple and known.

Similarly, businesses push on with "Big Data" for fear of missing
out, but without any real understanding of what they really
need. Solutions are commissioned that aspire to the heights of Facebook and Google
whilst in reality they fumble for a real business-value providing use-case.

Don't get me wrong, I'm 100% all for learning new paradigms, languages
and frameworks. This is just a reminder, as much to myself as anyone
else, to take a moment to think past the biases that may lead us to
make some regrettable, albeit well-intentioned and over-excited, choices.


## References

* [You Are Not So
  Smart](https://youarenotsosmart.com/2013/05/23/survivorship-bias/)
* [How Not to Be Wrong: The Hidden Maths of Everyday Life](https://www.amazon.co.uk/dp/071819604X)
* [High performance envy/web scale envy](https://www.thoughtworks.com/radar/techniques/high-performance-envy-web-scale-envy)
