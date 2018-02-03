# Vault Spec

When you use Vault as part of your application, your application usually
expects some kind of structure inside the secret store. Usually you'd have a
set of secrets each containing a handful of properties. But how do you prepare
this structure inside your Vault? This is where vault-spec comes in.

Using a single configuration file you can configure Vault to contain a certain
structure within the key-value store and also a set of policies that your app
will be using.

## Usage

1. Define a spec file (e.g. `sample.vaultspec.yaml`) for your Vault instance.
   You can find details on format for the configuration file down below.
2. Run `vault-spec validate -f sample.vaultspec.yaml` in order to make sure
   your specification file is valid.
3. Run `vault-spec apply -f sample.vaultspec.yaml` in order to apply the
   specification to the Vault server. Note that you will have to have
   `VAULT_ADDR` and `VAULT_TOKEN` set up with capabilities to create policies
   and write to the specified secret paths.


## Configuration format

The vault-spec configuration file is written in YAML and looks similar to what
you can find inside the `sample.vaultspec.yaml` file. Every file has to first
specify a version:

### `version`

Right now, only `"1"` is supported.

### `spec`

The only other top-level property of a specification is `spec`, which houses
`secrets`.

### `secrets`

This section allows you to define all the secrets and their respective
properties that you Vault should include according to the spec. If you already
know JSON Schema, this is pretty similar:

```
spec:
  secrets:
    - path: secrets/path/to/secret
      properties:
        username: ...
        password: ...
```

### `secrets[].path`

This is where you define where inside Vault's key-value store the secret should
be situated. This field is mandatory.


### `secrets[].properties`

Each secret consists of at least one property which you can define here using a
dictionary-like structure. Properties consist of a type (default `string`) and
also have a `default` value. The user is prompted to enter a value which is
then validated against this type. If the user's input should not be printed,
use `input: password`.

Summarizing, each property can have the following settings:

* `type`: The type of the property (e.g. `string`, `int`, `float`, ...),
  default: `string`)
* `input`: Defines how the user should be prompted for input (`password` or
  `default`).
* `help`: When the user is prompted for input, this message is shown alongside
  the property's name.
* `default`: The default value that should be used if the user simply hits
  [ENTER] when prompted.