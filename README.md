# Paketo Buildpack for Apache SkyWalking

## Buildpack ID: `paketo-buildpacks/apache-skywalking`
## Registry URLs: `docker.io/paketobuildpacks/native-image`

The Paketo Buildpack for Apache SkyWalking is a Cloud Native Buildpack that contributes the [Apache SkyWalking][s] Agent and configures it to connect to the service.

[s]: https://skywalking.apache.org

## Behavior

This buildpack will participate if all the following conditions are met

* A binding exists with `type` of `ApacheSkyWalking`

The buildpack will do the following for Java applications:

* Contributes a Java agent to a layer and configures `$JAVA_TOOL_OPTIONS` to use it
* Transforms the contents of the binding secret to system properties with the pattern `-Dskywalking.<KEY>=<VALUE>`

## Bindings

The buildpack optionally accepts the following bindings:

### Type: `dependency-mapping`

| Key                   | Value   | Description                                                                                       |
| --------------------- | ------- | ------------------------------------------------------------------------------------------------- |
| `<dependency-digest>` | `<uri>` | If needed, the buildpack will fetch the dependency with digest `<dependency-digest>` from `<uri>` |

## License

This buildpack is released under version 2.0 of the [Apache License][a].

[a]: http://www.apache.org/licenses/LICENSE-2.0
