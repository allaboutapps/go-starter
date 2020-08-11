# `/assets`

Other assets to go along with your repository (images, logos, etc).

https://github.com/golang-standards/project-layout/tree/master/assets

### all about apps specific requirement regarding `/assets/mnt`

`/assets/mnt` is special, as this is typically our PV rw mountpoint in our infrastucture.

**Do not check-in files there** (it's also `.gitignore`d)!

Instead, use this path as your default place to write user generated content that needs to be persisted. This folder will get **shadowed** by the PV mount in our infrastructure and will be available by all replicas through NFS.