# legacyflag – EXPERIMENTAL

Provides a pflag wrapper that addresses pain-points of backwards-compatible `ComponentConfig` migration.

## Development Tips

If you modify the codegen templates in `hack/gen/gen.go`, or update the 
generator config in `hack/gen/config.json` remember to run
`go run hack/gen/gen.go --config hack/gen/config.json` to regenerate source 
files and tests.

## Community, discussion, contribution, and support

Learn how to engage with the Kubernetes community on the [community page](http://kubernetes.io/community/).

You can reach the maintainers of this project at:

- [Slack channel](https://kubernetes.slack.com/messages/wg-component-standard)
- [Mailing list](https://groups.google.com/forum/#!forum/kubernetes-wg-component-standard)

### Code of conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](code-of-conduct.md).
