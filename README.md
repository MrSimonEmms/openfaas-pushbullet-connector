# OpenFaaS PushBullet Connector

Trigger OpenFaaS jobs from [Pushbullet](https://pushbullet.com)

## What is Pushbullet?

Pushbullet is a way of connecting your devices together. It has a simple sharing mechanism which allows you to push
between desktop, tablet and phone. It allows you to push notes, links and more.

From an OpenFaaS point of view, it's can be used as a managed message queue and allows you to send events from any
events with their simple and effective interface.

It's also free and without ads. I've used it for nearly a decade.

## Getting started

### Getting a Pushbullet token

1. Go to your [Pushbullet account homepage](https://www.pushbullet.com/#settings/account)
2. Hit the "Create Access Token" button
3. The access token will be displayed once

### Create a channel

Pushbullet allows you to create "channels". These can be subscribed to by anyone, but can only be pushed to by the
owner. This is perfect for triggering OpenFaaS events.

To create a channel, go to the [Pushbullet Create Channel](https://www.pushbullet.com/my-channel) page. This allows
you three options - from the point of view of this connector, the only important one is the `tag`.

- Tag: this **MUST** be set to the OpenFaaS function's topic name. This is globally-unique.
- Channel name: this is the human-readable channel name. Make it useful.
- Description: the description, obviously.

### Topic name

The [OpenFaaS Connector SDK](https://github.com/openfaas/connector-sdk) sends events to an OpenFaaS function with a
`topic` annotation.

The OpenFaaS Pushbullet Connector treats the Pushbullet `tag` and the OpenFaaS `topic` as the same. That means that
the event is sent to the OpenFaaS function with a `topic` annotation that is identical to the Pushbullet `tag`.

Pushbullet tags are globally-unique.

If the connector can't match the tag to a channel or that the channel doesn't match an OpenFaaS function, the event
is simply ignored.

## Kubernetes

See the [Helm chart](/chart/openfaas-pushbullet-connector) for details

## Todo

- [ ] Trigger a response when a function has finished
- [ ] Improved async OpenFaaS functions 

PRs welcome.
