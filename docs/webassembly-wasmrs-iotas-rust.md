# Extending services with WebAssembly

Nanobus is a way to stitch together logic, no matter where it exists or what form it's in. We call bundles of such logic "iotas" meaning "small things". Iotas are simply a collection of interfaces that expose actions (i.e. functions). They're just like language libraries or modules you've used before, except with an interface that fits distributed, cross-language, cross-platform usage.

Iotas need to communicate with known protocols and WebAssembly iotas use RSocket via wasmRS to establish bidirectional communication between the WebAssembly guest and the iota host, nanobus.

## Starting off

Use [apex] to start a new project for your language of choice. This command uses the template found in the path `templates/rust` from the repository `git@github.com:nanobus/iota.git`. Any git repository can be an apex template.

Use the project name `jinja` for this example. We'll be building a template-rendering iota.

```
apex new git@github.com:nanobus/iota.git -p templates/rust jinja
```

Apex is useful for more than just creating boilerplate projects. It's a code generation tool that revolves around the apex definition language written in `.axdl` files. Apex is a minimal IDL that lets you define your APIs, types, and documentation in one place and generate code, markup, documentation, boilerplate, integration, etc with the `apex` cli.

The template's definition file is `apex.axdl` and looks like this:

```apex
namespace "jinja"

interface MyModule @service {
  action1(first:string, last:string): string
}
```

This defines a module named `MyModule` with one action named `action1` that takes two strings and returns one string. The `@service` annotation tells the iota code generators which interface is a public interface vs imported or depended upon interfaces.

For this iota, we're building a template renderer using `jinja`-like syntax. Set our namespace to `jinja` (if it's not already) and our interface name to `Jinja`. Remove the sample action `action` and replace it with the following: `render(template:string, data:{string:string}): string`. That is, we're defining an action named `render` that takes two inputs, the template source as a string and the data for the template as a map of string:string pairs. The output of our action is a string.

Our apex should now look like this:

```apex
namespace "jinja"

interface Jinja @service {
  render(template:string, data:{string:string}): string
}
```

Run `just codegen` to generate the boilerplate for our project.

```console
$ just codegen
apex generate
INFO Writing file ./src/actions/jinja/render.rs (mode:644)
INFO Writing file ./src/lib.rs (mode:644)
INFO Writing file ./src/error.rs (mode:644)
INFO Writing file ./src/actions/mod.rs (mode:644)
INFO Formatting file ./src/actions/jinja/render.rs
INFO Formatting file ./src/error.rs
INFO Formatting file ./src/lib.rs
INFO Formatting file ./src/actions/mod.rs
```

The most important file to note here is `./src/actions/jinja/render.rs`. This is where our action will be implemented.

Before we implement anything, add a rust dependency to render templates: `minijinja`

```console
$ cargo add minijinja
    Updating crates.io index
      Adding minijinja v0.27.0 to dependencies.
             Features:
             + builtins
             + debug
[ features snipped...]
```

Now we write normal rust code to render a template we get as `input.template` with the data we get from `input.data`. Our output is a string.

```rs
use minijinja::Environment;

use crate::actions::jinja_service::render::*;

pub(crate) async fn task(input: Inputs) -> Result<Outputs, crate::Error> {
    let mut env = Environment::new();
    env.add_template("root", &input.template).unwrap();

    let template = env.get_template("root").unwrap();
    let rendered = template.render(input.data).unwrap();
    Ok(rendered)
}
```

Build our iota with the command:

```
$ just build
[cargo output snipped...]
$ ls build/
jinja.wasm
```

To test our iota we need to run it in a wasmrs-capable host. Nanobus natively supports wasmrs and we can set up a configuration that delegates to our wasm file.

Make a new file called `bus.yaml` and add the following:

```yaml
id: jinja
version: 0.0.1
main: build/jinja.wasm
```

Now we can run `nanobus invoke` and pipe our input through it to see our iota in action.

```console
$ echo '{"template":"hello {{ name }}!", "data":{"name": "world"}}' | nanobus invoke jinja.Jinja::render
"hello world!"
```

[apex]: https://apexlang.io
[just]: https://github.com/casey/just#packages
[starter]: https://github.com/nanobus/starter-template
