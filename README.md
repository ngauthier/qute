# Objects

## Poller (specific)

Queries db for builds to create jobs to give to dispatcher

* Outgoing build source channel to dispatcher
* Incoming build finished channel from dispatcher
* List of all builds running (map on uuid?)

1. Query for running builds not in list (initially empty) (that have not been ticked for some time amount, possibly based on state and prioritized)
1. Add build to list and send build outgoing to dispatcher
1. Receive build incoming as finished from dispatcher, remove from list (let the query pick it back up)

Handle panic by closing outgoing channel and waiting for incoming channel to close, and restart loop

Catch O.S. signals like ListenAndServe does so it can graceful shutdown (similar to the panic restart, but without the restart)

## Dispatcher (generic)

* Incoming job source channel
* Channel of Worker Channels
* Outgoing job finished channel

## Worker (generic)

* Channel of Worker Channels (put yourself when done)
* Channel for one build

1. Get job
1. Do job
1. Send job back
1. Put self back on worker channel channel

## Job (generic)

* Build wrapper w/ api
* Impls `Do() error`

## XJob (specific)

Typed Job that was instantiated by Poller based on build state


# Notes

* `select` can be used for non blocking send as well as receive
* `range` can be used to loop on channel receives
* http://www.slideshare.net/cloudflare/a-channel-compendium slide 38: load balancer has a fixed worker pool
