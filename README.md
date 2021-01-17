# Residential DNS

![Black Lives Matter](https://img.shields.io/badge/BLM-Black%20Lives%20Matter-black)
![CircleCI](https://img.shields.io/circleci/build/github/patoarvizu/residential-dns.svg?label=CircleCI) ![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/patoarvizu/residential-dns.svg) ![Docker Pulls](https://img.shields.io/docker/pulls/patoarvizu/residential-dns.svg) ![Keybase BTC](https://img.shields.io/keybase/btc/patoarvizu.svg) ![Keybase PGP](https://img.shields.io/keybase/pgp/patoarvizu.svg) ![GitHub](https://img.shields.io/github/license/patoarvizu/residential-dns.svg)

<!-- TOC -->

- [Residential DNS](#residential-dns)
  - [Intro](#intro)
    - [WARNING](#warning)
  - [Arguments](#arguments)
  - [Credentials](#credentials)
  - [Multi-architecture images](#multi-architecture-images)

<!-- /TOC -->

## Intro

This is a simple controller, that discovers the public IP address of the node where it's running from (using the [ipify API](https://www.ipify.org/)), and creates a Route 53 `A` record pointing to that public IP address.

The value of the controller is to create a consistent hostname that can be used in cases where acquiring a fixed IP address is not practical.

Although this can be run as a one-off Go script, the intended design is to run it as a Docker container (optimally, in Kubernetes or other container orchestration platform) and it will continuously keep a record updated with the value of the public IP address of your home router, which may change after a reboot.

### WARNING

Exposing your home network to the internet will always be a risk. Only use this if you're aware of what you're doing and you have taken the appropriate steps to mitigate the risk, like protecting your home network with a firewall, or routing the traffic through a [CDN](https://en.wikipedia.org/wiki/Content_delivery_network) (Content delivery network).

## Arguments

Command-line argument  | Default | Description
-----------------------|---------|------------
`-hosted-zone-id`      |         | Route 53 hosted zone id
`-record-name`         |         | DNS record name
`-ttl`                 |      60 | DNS record TTL
`-sync-period-minutes` |      15 | The amount of time, in minutes, to wait between syncs

## Credentials

The controller will need AWS credentials and will discover them using the [default credential provider chain](https://docs.aws.amazon.com/sdk-for-java/v1/developer-guide/credentials.html), so they can either use static credentials provisioned manually, or another dynamic credential injection mechanism, such as Vault's [AWS Secrets Engine](https://www.vaultproject.io/docs/secrets/aws), [kube2iam](https://github.com/jtblin/kube2iam), [kiam](https://github.com/uswitch/kiam), etc., depending on the runtime.

The only permission required is `route53:ChangeResourceRecordSets`.

## Multi-architecture images

Manifests published with the semver tag (e.g. `patoarvizu/residential-dns:v0.0.0`), as well as `latest` are multi-architecture manifest lists. In addition to those, there are architecture-specific tags that correspond to an image manifest directly, tagged with the corresponding architecture as a suffix, e.g. `v0.0.0-amd64`. Both types (image manifests or manifest lists) are signed with Notary as described above.

Here's the list of architectures the images are being built for, and their corresponding suffixes for images:

- `linux/amd64`, `-amd64`
- `linux/arm64`, `-arm64`
- `linux/arm/v7`, `-arm7`