# Residential DNS

![Black Lives Matter](https://img.shields.io/badge/BLM-Black%20Lives%20Matter-black)
![CircleCI](https://img.shields.io/circleci/build/github/patoarvizu/residential-dns.svg?label=CircleCI) ![GitHub tag (latest SemVer)](https://img.shields.io/github/tag/patoarvizu/residential-dns.svg) ![Docker Pulls](https://img.shields.io/docker/pulls/patoarvizu/residential-dns.svg) ![Keybase BTC](https://img.shields.io/keybase/btc/patoarvizu.svg) ![Keybase PGP](https://img.shields.io/keybase/pgp/patoarvizu.svg) ![GitHub](https://img.shields.io/github/license/patoarvizu/residential-dns.svg)

<!-- TOC -->

- [Residential DNS](#residential-dns)
  - [Intro](#intro)
    - [WARNING](#warning)
  - [Arguments](#arguments)
  - [Credentials](#credentials)
  - [For security nerds](#for-security-nerds)
    - [Docker images are signed and published to Docker Hub's Notary server](#docker-images-are-signed-and-published-to-docker-hubs-notary-server)
    - [Docker images are labeled with Git and GPG metadata](#docker-images-are-labeled-with-git-and-gpg-metadata)
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

## For security nerds

### Docker images are signed and published to Docker Hub's Notary server

The [Notary](https://github.com/theupdateframework/notary) project is a CNCF incubating project that aims to provide trust and security to software distribution. Docker Hub runs a Notary server at https://notary.docker.io for the repositories it hosts.

[Docker Content Trust](https://docs.docker.com/engine/security/trust/content_trust/) is the mechanism used to verify digital signatures and enforce security by adding a validating layer.

You can inspect the signed tags for this project by doing `docker trust inspect --pretty docker.io/patoarvizu/residential-dns`, or (if you already have `notary` installed) `notary -d ~/.docker/trust/ -s https://notary.docker.io list docker.io/patoarvizu/residential-dns`.

If you run `docker pull` with `DOCKER_CONTENT_TRUST=1`, the Docker client will only pull images that come from registries that have a Notary server attached (like Docker Hub).

### Docker images are labeled with Git and GPG metadata

In addition to the digital validation done by Docker on the image itself, you can do your own human validation by making sure the image's content matches the Git commit information (including tags if there are any) and that the GPG signature on the commit matches the key on the commit on github.com.

For example, if you run `docker pull patoarvizu/residential-dns:f61593ce4a5bd89936742a87f4e70d0cdba7d3d7` to pull the image tagged with that commit id, then run `docker inspect patoarvizu/residential-dns:f61593ce4a5bd89936742a87f4e70d0cdba7d3d7 | jq -r '.[0].Config.Labels'` (assuming you have [jq](https://stedolan.github.io/jq/) installed) you should see that the `GIT_COMMIT` label matches the tag on the image. Furthermore, if you go to https://github.com/patoarvizu/residential-dns/commit/f61593ce4a5bd89936742a87f4e70d0cdba7d3d7 (notice the matching commit id), and click on the **Verified** button, you should be able to confirm that the GPG key ID used to match this commit matches the value of the `SIGNATURE_KEY` label, and that the key belongs to the `AUTHOR_EMAIL` label. When an image belongs to a commit that was tagged, it'll also include a `GIT_TAG` label, to further validate that the image matches the content.

Keep in mind that this isn't tamper-proof. A malicious actor with access to publish images can create one with malicious content but with values for the labels matching those of a valid commit id. However, when combined with Docker Content Trust, the certainty of using a legitimate image is increased because the chances of a bad actor having both the credentials for publishing images, as well as Notary signing credentials are significantly lower and even in that scenario, compromised signing keys can be revoked or rotated.

Here's the list of included Docker labels:

- `AUTHOR_EMAIL`
- `COMMIT_TIMESTAMP`
- `GIT_COMMIT`
- `GIT_TAG`
- `SIGNATURE_KEY`

## Multi-architecture images

Manifests published with the semver tag (e.g. `patoarvizu/residential-dns:v0.0.0`), as well as `latest` are multi-architecture manifest lists. In addition to those, there are architecture-specific tags that correspond to an image manifest directly, tagged with the corresponding architecture as a suffix, e.g. `v0.0.0-amd64`. Both types (image manifests or manifest lists) are signed with Notary as described above.

Here's the list of architectures the images are being built for, and their corresponding suffixes for images:

- `linux/amd64`, `-amd64`
- `linux/arm64`, `-arm64`
- `linux/arm/v7`, `-arm7`