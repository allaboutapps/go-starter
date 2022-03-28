# `/cmd`

Main applications for this project.

Don't put a lot of code in the application directory. If you think the code can be imported and used in other projects, then it should live in the `/pkg` directory. If the code is not reusable or if you don't want others to reuse it, put that code in the `/internal` directory. You'll be surprised what others will do, so be explicit about your intentions!

We manage our applications via cobra ([`cobra-cli`](https://github.com/spf13/cobra-cli) is installed within the `Dockerfile` within the development stage, `/cmd` consumes the core [`cobra`](https://github.com/spf13/cobra) library), see:
* https://github.com/spf13/cobra#getting-started
* https://github.com/spf13/cobra-cli/blob/main/README.md#add-commands-to-a-project

Also see https://github.com/golang-standards/project-layout/tree/master/cmd