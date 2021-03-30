Gcore cloud API client
====================================

Command line client to GCore cloud API.

Installation
------------------------------------

Clone repo and run `make build`.

Getting started
------------------------------------

You will need to set the following env:
```bash
export GCLOUD_USERNAME=username
export GCLOUD_PASSWORD=secret
export GCLOUD_PROJECT=1
export GCLOUD_REGION=1
export GCLOUD_AUTH_URL=https://api.gcdn.co
export GCLOUD_API_URL=https://api.cloud.gcorelabs.com
export GCLOUD_CLIENT_TYPE=platform
```

* **GCLOUD_USERNAME** - username
* **GCLOUD_PASSWORD** - user's password
* **GCLOUD_PROJECT** - project id
* **GCLOUD_REGION** - region id
* **GCLOUD_AUTH_URL** - authentication url, you could use the same as in example above
* **GCLOUD_API_URL** - api url, you could use the same as in example above
* **GCLOUD_CLIENT_TYPE** - client type, you could use the same as in example above

After setting the env, use `-h` key to retrieve all available commands:
```bash
./gcoreclient -h

NAME:
   gcoreclient - GCloud API client

   Environment variables example:

   GCLOUD_AUTH_URL=
   GCLOUD_API_URL=
   GCLOUD_API_VERSION=v1
   GCLOUD_USERNAME=
   GCLOUD_PASSWORD=
   GCLOUD_REGION=
   GCLOUD_PROJECT=


USAGE:
   gcoreclient [global options] command [command options] [arguments...]

VERSION:
   v0.2.11

COMMANDS:
   network        GCloud networks API
   task           GCloud tasks API
   keypair        GCloud keypairs V2 API
   volume         GCloud volumes API
   subnet         GCloud subnets API
   flavor         GCloud flavors API
   loadbalancer   GCloud loadbalancers API
   instance       GCloud instances API
   heat           Gcloud Heat API
   securitygroup  GCloud security groups API
   floatingip     GCloud floating ips API
   port           GCloud ports API
   snapshot       GCloud snapshots API
   image          GCloud images API
   region         GCloud regions API
   project        GCloud projects API
   keystone       GCloud keystones API
   quota          GCloud quotas API
   limit          GCloud limits API
   cluster        Gcloud k8s cluster commands
   pool           Gcloud K8s pool commands
   help, h        Shows a list of commands or help for one command
```
