# `/scripts`

Scripts to perform various build, install, analysis, etc operations.

These scripts keep the root level Makefile small and simple.

https://github.com/golang-standards/project-layout/tree/master/scripts

Examples:

* https://github.com/kubernetes/helm/tree/master/scripts
* https://github.com/cockroachdb/cockroach/tree/master/scripts
* https://github.com/hashicorp/terraform/tree/master/scripts

Please note that this scripts are not available in a final product. Use `/cmd` instead if you need to execute your script an a live environment.