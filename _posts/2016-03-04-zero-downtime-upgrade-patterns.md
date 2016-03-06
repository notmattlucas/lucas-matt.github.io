---
layout: post
title: Zero Downtime Upgrade Problems and Patterns
description: "A survey of the problems that various technologies face in this regard, and detail common patterns for tackling them."
modified: 2016-03-04
tags: [upgrade, docker, linux, production]
image:
  feature: downtime.jpg
---
In the modern software landscape it is increasingly important, not to mention expected, that services be upgraded without causing downtime to the end user. In an industry such as IPTV to make the service unavailable to the end user for a period of time, even in the middle of the night, is generally unacceptable.

This page will survey the problems that various technologies face in this regard, and detail common patterns for tackling them.

This page will survey the problems that various technologies face in this regard, and detail common patterns for tackling them.

It's worth noting here that most of these approaches are essentially just variations on the Blue Green Deployment strategy, as described by Martin Fowler [here](http://martinfowler.com/bliki/BlueGreenDeployment.html).

##Stateless Systems / Considering Only the Application Layer

A stateless system, that is one without any kind of persistent store, is obviously the simplest case we can consider for a zero downtime upgrade and will cover many principles that can be used in all other cases, when considering the application in isolation.

###Remove-Upgrade-Add

The most obvious strategy, and the one that is probably most widely used, is that of removing the service from a load balancer, upgrading that service, and once upgraded re-enabling it. Of course this assumes you have some redundancy in your platform to enable you to take parts of it down whilst still serving traffic.

{:refdef: style="text-align: center;"}
![Remove Add Upgrade Animation]({{ site.url }}/images/remove-upgrade-add.gif){: .center-image }
{: refdef}

Unfortunately, this is a pretty coarse approach to upgrade, and is generally achieved manually by some poor engineer in the middle of the night. As manual is usually synonymous with bug-prone then automating this would be desirable, but isn't trivial in a traditional bare-metal deployment as managing running services and state over potentially many nodes is quite tricky and bug-prone in itself.

###Fixing the Restart Problem

####Master - Workers Model

World class http servers (e.g. nginx, unicorn) and process managers (e.g. node's pm2) have inbuilt features that allow for zero-downtime upgrades within the context of a single machine, avoiding the tedium of the above remove-upgrade-add cycle but achieving the same end goal (remember we're still talking application layer only, so no state to worry about).

The core architecture that allows them to achieve this is a master-workers model. In this model we have one core process performing all privileged operations (such as reading configuration or binding to ports) and a number of workers and helper processes (the app itself). The point here, is that this allows for the workers to be switched out behind the scenes, the master effectively working as a software load balancer.

> [Inside Nginx - how we designed for performance/scale](https://www.nginx.com/blog/inside-nginx-how-we-designed-for-performance-scale/)

![Nginx on the fly upgrade]({{ site.url }}/images/nginx-upgrade.png)

> "Nginx’s binary upgrade process achieves the holy grail of high-availability – you can upgrade the software on the fly, without any dropped connections, downtime, or interruption in service.

> The binary upgrade process is similar in approach to the graceful reload of configuration. A new NGINX master process runs in parallel with the original master process, and they share the listening sockets. Both processes are active, and their respective worker processes handle traffic. You can then signal the old master and its workers to gracefully exit."

There are tools that are able to achieve this natively for your own project (the aforementioned pm2) and also as a generic solution across various platforms ([socketmaster](https://github.com/zimbatm/socketmaster), [einhorn](https://github.com/stripe/einhorn)).

You could also achieve the equivalent of this architecture by parallel running both versions of your app on the same node, and using <strong>nginx reload</strong> to seamlessly switch your configuration between the two.

###so_reuseport

Introduced in Linux kernel 3.9 was the SO_REUSEPORT option. As long as the first server sets this option before binding its socket, then any number of other servers can also bind to the same port if they also set the option beforehand - basically no more "Address already in use" errors. Once enabled the SO_REUSEPORT implementation will distribute connections evenly across all of the processes.

{:refdef: style="text-align: center;"}
<img src="{{site.url}}/images/so_reuseport.png" width="500" />
{: refdef}

Of course, the implication here is that we can simply spin up multiple versions of our application, an old and new; drain the old of requests and once drained shut it down.

##Stateful Systems / Evolving the Database

It's when we get to stateful systems that we begin to hit problems with zero downtime upgrades, or live migration. There is much to consider here, as we have various data-store properties and numerous different uses cases, especially in a micro-service model where each service is tuned in a very specific way. Inevitably most of the available patterns do not fit all circumstances.

First let's consider the various dimensions along which this problem can vary:

**Mutability of Data**

* Read only systems - *metadata retrieval, or notification systems*
* Write only systems - *immutable data stores, event sourcing systems*
* Read-Write systems - *account management, bookmarking, etc*

**Data Store Properties**

* Schema based DB system - *Oracle, Cassandra*
* Schema-less - *MongoDB, Redis*
* Archive (essentially non-interactive) data store - *e.g. HDFS with Avro*

**Evolution Type**

* Addition of a field
* Removal or rename of a field
* Removal of a record
* Refactoring of a record - *e.g. splitting into multiple*

###Polymorphic Schema / Lazy Migration

A schema-less database, such as MongoDB, can support a pattern known as a "Polymorphic Schema", also known as a "single-table model". Rather than all documents in a collection being identical, and so adhering to an implicit schema, the documents can be similar, but not identically structured. These schema, of course, map very closely to object-oriented programming principles. This style of schema feeds well into a schema evolution plan.

{:refdef: style="text-align: center;"}
![MongoDB Applied Design Patterns]({{site.url}}/images/MongoDB-Applied-Design-Patterns.png)
{: refdef}

A traditional RDBMS will evolve its schema through a set of migration scripts. When the system is live, these migrations can become complex or require database locking operations, resulting in periods of extended downtime. This model is also possible in Mongo DB, and similar, but we have alternatives that may suit us better.

By using a polymorphic schema we can account for the absence of a new field through defaults. Similar behavior can be implemented for field removal also, and these kinds of transformations are handled elegantly by Mongo ORM libraries if you want to avoid coding this logic yourself.

By applying these transformations lazily, as records are naturally read and written back to the database, the schema is slowly transformed into the latest. By running through all records in a collection, we would be able to force a full upgrade of the full collections.

Of course there are other, more specific, challenges depending on the type of transformation required:

* How do we perform collection refactoring exercises (e.g. splitting out a single collection into multiple related collections) - *would a _version field on an object allow us to trigger migrations more effectively in this case?*

* How would we perform queries based upon that new field? - *by forcing a read-write of all records, and therefore completing the transformation whilst the system was still alive?*

### Forward Only Migrations

For database systems that enforce a schema, but a principle that is equally as applicable to all persistence stores, is the concept of forward only migrations.

{:refdef: style="text-align: center;"}
<img src="{{site.url}}/images/forward-migration.jpg" width="500" />
{: refdef}

**Never roll back == No messy recovery**

In this practice every database migration performed should be compatible with both the **new** version of the code *and* the **previous** one. Then if you have to roll back a code release the previous version is perfectly happy using the new version of the software.

As you can imagine, this will require some strict adherence to convention. For example, dropping a column:

**Dropping a Column**

* Release a version of the code that doesn't use the column.
* Ensure that it is stable, and won't be rolled back.
* Do a second release that has a migration to remove the column now that it isn't being used.

**Rename a Column**

* Release a version that adds a new column with the new name, changing the code to write to both but read from the old.
* Employ a batch job to copy data from the old column to the new column - one row at a time to avoid locking the database.
* Release a version that reads from and writes to the new column.
* In the next release, drop the column

This technique seems most successfully employed in a continuous delivery environment ([I](http://engineering.skybettingandgaming.com/2016/02/02/how-we-release-so-frequently/))

###Expansion / Upgrade / Contraction Scripts

This method is similar to the last, but it's organized slightly differently. Here we use two different sets of migration scripts:

**Expansion Scripts**

* Apply changes to the documents safely, that **do not** break backwards compatibility with the existing version of the application.
* e.g. Adding, copying, merging splitting fields in a document.

**Contraction Scripts**

* Clean up any database schema that is not needed after the upgrade.
* e.g. removing fields from a document

{:refdef: style="text-align: center;"}
<img src="{{site.url}}/images/expand-contract.png" width="500" />
{:refdef}

The process we perform here is:

1. **Expand** - Run expansion scripts before upgrading the application
2. **Upgrade** - Upgrade the cluster one node at a time
3. **Contract** - Once the system has been completely upgraded and deemed stable. Typically, contractions can be run, say days/weeks after complete validation

Again, this does not require any DB rollback. The reason for attempting this is that reversing DB changes can lead to potential loss of data or leaving your system in an inconsistent state. It is safer to rollback the application without needing to rollback DB changes as expansions are backwards compatible.

###One Big Caveat for the Above

Most of the above solutions require that upgrades are nice and discrete. That is, they move from version to version in an incremental way, not skipping any releases so that we can control schema compatibility in a pair-wise fashion.

{:refdef: style="text-align: center;"}
<img src="{{site.url}}/images/incremental.png" width="400" />
{: refdef}

This is all fine for a continuous delivery, or cloud-based shop, where upgrades are obviously linear and wouldn't be applied any other way. However what if you have a user-base of several parallel deployments all leap-frogging over various versions, as below.

{:refdef: style="text-align: center;"}
<img src="{{site.url}}/images/good_skip.png" width="400" />
{: refdef}

One solution is simply to enforce your upgrade path more clearly. So, for example, any app-only change would be a simple revision bump. But any schema change that needs to be handled by one of the above methods would be a larger bump.

It would have to be mandatory to hit every one of these larger version numbers to preserve this incremental schema-change model.

{:refdef: style="text-align: center;"}
<img src="{{site.url}}/images/bad_skip.png" width="400" />
{: refdef}

##Immutable Infrastructure - Docker and Friends

In the world of containers, as deployments are immutable, then a number of the above techniques are not applicable. But the same blue-green principles apply.

{:refdef: style="text-align: center;"}
<img src="{{site.url}}/images/docker_zdu.png" width="500" />
{: refdef}

[HAProxy](http://www.haproxy.org/) is a web proxy and load balancer that allows you to change the configuration, restart the process, and (most importantly) have the old process die automatically when there is no more traffic running through it. So a simple upgrade would go as follows:

* Start a new version of the container
* Tell HAProxy to direct all new traffic to the new container whilst maintaining existing connections to the previous.
* All new users get access to the new code immediately, whilst existing users are not cut-off from their in-progress requests.

Of course, this can all be automated, but that's for another time.

##References

* [No Downtime Database Migrations](http://www.planetcassandra.org/blog/no-downtime-database-migrations/)
* [Blue Green Deployment](http://martinfowler.com/bliki/BlueGreenDeployment.html)
* [Implementing Blue-Green Deployments with AWS](https://www.thoughtworks.com/insights/blog/implementing-blue-green-deployments-aws)
* [Zero Downtime Clojure Deployments](http://www.uswitch.com/tech/zero-downtime-clojure-deployments/)
* [Zero Downtime Deployments (in the context of Docker)](https://docs.quay.io/solution/zero-downtime-deployments.html)
* [Kubernetes Rolling Update Example](http://kubernetes.io/v1.1/docs/user-guide/update-demo/README.html)
* [Framework Agnostic Zero Downtime Deployment with Nginx](http://openmymind.net/Framework-Agnostic-Zero-Downtime-Deployment-With-Nginx/)
* [Inside NGINX: How We Designed for Performance & Scale](https://www.nginx.com/blog/inside-nginx-how-we-designed-for-performance-scale/)
* [Zero Downtime, Instance Deployment and Rollback](http://www.ebaytechblog.com/2013/11/21/zero-downtime-instant-deployment-and-rollback/)
* [The SO_REUSEPORT Socket Option](https://lwn.net/Articles/542629/)
* [How We Release So Frequently](http://engineering.skybettingandgaming.com/2016/02/02/how-we-release-so-frequently/)
* [Continuous Delivery and the Database](https://www.simple-talk.com/sql/database-administration/continuous-delivery-and-the-database/)
* [Midas On the Fly Migration Tool for Mongo DB](http://www.slideshare.net/DhavalDalal/midas-onthefly-schema-migration-tool-for-mongodb)
* [Akka Persistence - Schema Evolution](http://doc.akka.io/docs/akka/snapshot/scala/persistence-schema-evolution.html)
* [MongoDB Applied Design Patterns [Book]](http://shop.oreilly.com/product/0636920027041.do)
* [Migrating Your Database Upon Deployments](https://blog.justjuzzy.com/2012/11/migrate-your-database-upon-deployments/)
