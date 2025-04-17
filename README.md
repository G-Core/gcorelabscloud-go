Gcore cloud API client
====================================

Command line client to Gcore cloud API.

Installation
------------------------------------

- Clone repo
- Add $GOPATH/bin into $PATH
- Run `make install`.

Getting started
------------------------------------

You will need to set the following env:
```bash
export GCLOUD_USERNAME=username
export GCLOUD_PASSWORD=secret
export GCLOUD_PROJECT=1
export GCLOUD_REGION=1
export GCLOUD_AUTH_URL=https://api.gcore.com/iam
export GCLOUD_API_URL=https://api.gcore.com/cloud
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
   v0.3.45-3-gde7ea60

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
   l7policy       GCloud l7policy API
   router         GCloud router API
   fixed_ip       GCloud reserved fixed ip API
   help, h        Shows a list of commands or help for one command

Contributing
------------------------------------

We welcome contributions to the gcorelabscloud-go project! Whether it's bug fixes, feature additions, or documentation improvements, your help is appreciated.

Please read our [Contributing Guidelines](CONTRIBUTING.md) for details on:
- How to submit bug reports
- How to submit feature requests
- How to submit pull requests
- Development setup and guidelines
- Code style and testing requirements

For any questions or support, please check our [documentation](https://gcore.com/docs/cloud) or contact [GCore Support](https://gcore.com/support).

License
------------------------------------

This project is licensed under the Mozilla Public License Version 2.0 - see the [LICENSE](LICENSE) file for details.
