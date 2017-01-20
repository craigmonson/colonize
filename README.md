# Colonize #

Colonize is a configurable, albeit opinionated way to organize and manage
your terraform templates.  It revolves around the idea of environments, and
allows you to organize templates, and template data around that common idiom.

Once it's been configured, it allows for hierarical templates and variables, and
the ability to organize them in a defined manageable way.

## Install ##

### setup via go ###
### add .colonize.yaml to root of project ###

## Quick Concepts ##

Colonize is opinionated about the structure of a project, so, for
this to be used you must structure your terraform project in a specific way.

#### Project Root ####

The project root, according to colonize, is whereever the `.colonize.yaml` file
is located.  This is typically in the root of your git repo, but doesn't have
to be.  Configuration will not be read any higher in the tree than the project
root.

#### .colonize.yaml ####

There is one, global `.colonize.yaml` file that configures _how_ clonize runs.
This should be located in the root of your project.  Colonize will walk up the
tree until it finds thi file, at which pint it'll assume that's the project
root.  Environment configuration will be read through the branches up to the
root.

#### Leaves vs Branches ####

A Leaf is an endpoint in the tree, and should contain all your _functional_
terraform code.  This is where you keep the templates that you craft in order
to manage infrastructure.  This should be familiar to you if you've used
terraform before.

A Branch is just a pathway to a leaf.  These contain _NO_ functional terraform
templates, and _ONLY_ contain configuration in the `env` directory.

Colonize will utilize all of the configuration in each `env` directory, in each
branch, from the root, to the leaf, when runnin colonize commands.  This would
allow you to configure, say, an account variable in the root, and that account
variable would be available to every template in the tree, without you having
to cut and paste it into every template.

#### env directories ####

Each branch in the directory tree, from the root to the leaf can have environment
specific files, all stored in the 'env' directory at that particular level.
These files will be combined and used at each point in the leaves, when running
terraform commands.  Colonize and Terraform will work together to utilize these
combined files when the terraform commands are executed.

#### combined / generated files ####

Colonize utilizes terraform, and the way terraform runs commands to provide
environment specific configuration fo your functional templates.  It does this
by combining files through the tree, and placing them in the workin directory
of the currently executing template.  Those files are all prefixed with an
underscore ("\_"), for example: `_combined_variables.tf`.

#### derived variables/values ####

Colonize creates several variable, and assigns them values, automatically from
the generated config of the project / tree.  You can also define your own
derived variables to be used in your templates as well.  Unlike terraform,
colonize allows you to create variables and values from already existing values.


## Configurable Files ##

Colonize expects several specially named files in the `env` directory.  Each
one allows you to configure your templates at runtime in any way you see fit.
These files are:
  * \<environment\>.tfvars
  * derived.tfvars
  * \*.tf
  * remote\_setup.sh

In addition, you can name the templates in your templte directory to have
environment specific templates as well:
  * foo.tf.dev
  * foo.tf.base

Lets look a little deeper into each of these files.

### \<environment\>.tfvars ###

These are the driving environment spcific variable assignment files that will
distinguish settings between your different environments in terraform.  Each
environment is expected to map directly to a specific tfvars file.  I think it's
best described through an example, so given the tree:

```bash
test
└── web
    ├── env
    │   ├── dev.tfvars
    │   ├── prod.tfvars
    │   ├── qa.tfvars
    ├── main.tf
```

Here our `web` is setup with one `main.tf` file, where it's assumed that we're
using terraforms variables, and variable substitutions to create a more modular
and reusable template.  Lets assume that web is spinning up a single instance,
and we've got our instance size set to a variable:

```
# main.tf
variable "size" {}
resource "fake_instance", "fake" {
  size = "${var.size}"
  ...
}
```

We can now specify our instance size in each of the environment specific files:

```
# dev.tfvars
size = "small"

# qa.tfvars
size = "medium"

# prod.tfvars
size = "large"
```

Now, when we simply run one of our colonize commands with the environment set
appropriately, colonize will set things up so that terraform uses the right
tfvars file: `colonize prep -environment=dev` would use the `dev.tfvars` file
when doing terraform things.

Colonize will store all of those variable assignments in the `_combined.tfvars`
file in the leaf.  It will _also_ generate a variables file for you, also in
the root, named `_combined_variables.tf`

**NOTE**: Colonize can only use single string assignment variables at the
moment.  **NO MAPS OR LISTS**

### derived.tfvars ###

To Aid in configuration, colonize allows for single pass derived variables,
meaning that colonize will pass over the derived variables once for substitution.
Ths allows you to create more varied values based off of already defined
variables; Something that terraform currently doesn't do.  Colonize will first
generate the environment files as noted above.  It will then combine the
derived files, then make substitutions utilizing the combined variable file
from above.  For simplicity sake, it uses the same syntax that terraform does
for variable interpolations, but it does **NOT** allow the use of anything but
the `${var.variable_name}`, so **NO FUNCTIONS**.  Lets take a look.
Given the tree:


```bash
test
└── web
    ├── env
    │   ├── dev.tfvars
    │   ├── prod.tfvars
    │   ├── qa.tfvars
    │   ├── derived.tfvars
    ├── main.tf
```

As in the example above, the different files have different values for the
variable `size`.  Lets utilize that, and the `environment` varible that
Colonize automatically generates for us.  So our derived file looks like:

```bash
tag_name = "web-${var.environment}-${var.size}"
```

So, we'll run `colonize prep -environment=dev`, and colonize will build both the
`_combined_derived.tfvars` and `_combined_derived.tf` files for you, with the
contents:
```
# _combined_derived.tfvars
tag_name = "web-dev-small"

# _combined_derived.tf
variable "tag_name" {}
```

We can now utilize this variable in our templates as normal.  Do note that in
many cases it's possible to just utilize terraform for variable interpolations,
but in some cases it might be beneficial to allow for derived variables to
simplify the templates.

### \*.tf files ###

Any tf file in the configuration directory `env` will be combined and placed
into the `_combined.tf` file in the leaf.  There is nothing fancy that happens,
it just combines all the .tf files it finds in the tree between the root and
leaf.

### remote\_setup.sh ###

Colonize looks for only one remote\_setup.sh file in the roots config (`env`)
directory.  Like the derived tfvars files above, this one also allows for
variable interpolation.  Colonize will read in, do the variable substitutions,
and write it out in the leaf directory to a file named `_remote_setup.sh`
The remote is used when colonize does it's thing with terraform, like plan,
apply and destroy.

### Leaf specific files ###

In the leaves, you can have distinct templates per environment if you need.  By
naming the files with the appropriate extensions, colonize will know which ones
to combine when it prepares for the run.

Files named in the pattern: `template_name.tf.<environment>` will be included
when the environment matches.  So, a file named `foo.tf.dev` would be included
in the terraform run only if then environment is set to dev.  If the
environment is set to anything else, then it won't be included.  If there is a
template that's named with the `.default` extension, then environments that
do *NOT* have a specific template, will use the default one.  These files
will be combined into the `_combined.tf` file.  Lets take a look at an example:


```bash
test
└── web
    ├── env
    │   ├── ...
    ├── main.tf
    ├── creds.tf.base
    ├── creds.tf.prod
    ├── db.tf.dev
```

In the example above, we've got 4 terraform templates.  When we run the command:
`colonize prep -environment=prod`, then Colonize will include `creds.tf.prod`,
and *ONLY* `creds.tf.prod` into the `_combined.tf` file.

If we, instead, run the command: `colonize prep -environment=dev`, then Colonize
will include `creds.tf.base` AND `db.tf.dev` into the `_combined.tf` file.

Since terraform uses any files named with the `.tf` extension, when terraform
commands are executed, it will use both `main.tf`, and `_combined.tf`.

#### 4.  TBD: There must be a make\_order.txt file in order to use the branch make ####

## Execution ##

The best thing to do when refering to the execution of commnds in colonize is
to review the inline documentation.  However, what follows is a quick 
overview on what's available via the colonize commands.

The commands _mostly_ build upon themselves, so follow this order:
  * init
  * prep
  * plan
  * apply
  * destroy
  * clean

**NOTE**: prep will be run automatically for the plan command.  This is to
allow for the closest similarities to the actual terraform commands, from which
colonize tries to mimic. (plan, apply, destroy)

### ```colonize init``` ###

The **init** command runs an interactive process to help initialize your Colonize
project. It will ask a series of questions and provide defaults for building your
`.colonize.yaml` file.

The following is output is an example of the Colonize init command's interactive 
process. Each variable is provided with a default value, where entering nothing
will result in accepting that variable's default.

```
Enter 'environments_dir' [env]:
Enter 'base_environment_ext' [default]:
Enter 'combined_vals_file' [_combined.tfvars]:
Enter 'combined_vars_file' [_combined_variables.tf]:
Enter 'combined_derived_vals_file' [_combined_derived.tfvars]:
Enter 'combined_derived_vars_file' [_combined_derived.tf]:
Enter 'combined_tf_file' [_combined.tf]:
Enter 'combined_remote_config_file' [_remote_setup.sh]:
Enter 'remote_config_file' [remote_setup.sh]:
Enter 'derived_file' [derived.tfvars]:
Enter 'vals_file_env_post_string' [.tfvars]:
```

After completing each variable, the init command will display each setting and
prompt you for acceptance of the settings. After the settings have been accepted,
a `.colonize.yaml` file will be created in the current directory, as well as the
selected `environments_dir`.


**Optional Command Flags**

`--accept-defaults`: This will accept all default values, automatically.


Example:

Running `colonize init --accept-defaults`, would result in the following directory
structure:

    .
    ├── .colonize.yaml
    └── env

The contents of the `.colonize.yaml` file would be as follows:

    ## Generated by Colonize init
    ---
    environments_dir: env
    base_environment_ext: default
    autogenerate_comment: This file generated by colonize
    combined_vals_file: _combined.tfvars
    combined_vars_file: _combined_variables.tf
    combined_derived_vals_file: _combined_derived.tfvars
    combined_derived_vars_file: _combined_derived.tf
    combined_tf_file: _combined.tf
    combined_remote_config_file: _remote_setup.sh
    remote_config_file: remote_setup.sh
    derived_file: derived.tfvars
    vals_file_env_post_string: .tfvars


Of course, if you run the itneractive process and make modifications to any of
the variable defaults, the `.colonize.yaml` file would match those settings 
that you selected.

### ```colonize prep --environment=<env>``` ###

The **prep** command is the workhorse of the colonize command.  It does all of
the combining and tree walking to generate files that the installed terraform
will utilize in it's plan / apply / destroy runs.  As one would expect, this
prepares terraform for the given environment ```<env>```

All of the generated files are prepended with the underscore ("\_"), so should
be easily identifiable upon completion of the execution.

It should be noted that once prep has been successfully executed, you should
be able to execute _any_ terraform command, and the generated files will be
utilized as expected.  Since colonize only runs a subset of the terraform
commands, you can execute **prep** and run any terraform command to execute
outside of colonize.  Since it's terraform that does the state data syncing,
everything should stay ok, but you should be _very_ careful with this approach,
as remote file setup etc may need to be manually handled.

Prep does 2 things via terraform:
  * It removes any ```.terraform``` directory (remote state) as the *first* step of the execution.
  * It executes ```terraform get -update``` as the *last* step of the execution.

The ```get -update``` isn't such a big deal, but it's **VERY** important to
note that **prep** will remove the ```.terraform``` directory, as, depending on
what non-colonize commads you've been executing, you may accidentally remove
non-sync'd state data.

### ```colonize plan --environment=<env>``` ###

**plan** wraps ```terraform plan```.  It's important to understand that
plan will execute **prep** first, regardless if prep has already been
run.  This is important to know, because prep will **delete the .terraform
direcory** as a first step.

### ```colonize apply``` ###

**apply** wraps ```terraform apply```, and runs it against the existing plan
that was created in the ```plan``` step.  So yes, in order to run ```apply```,
you need to run ```plan``` first.

### ```colonize destroy --environment=<env>``` ###

**destroy** wraps ```terraform destroy```, and will fully destroy the template
stack.  It does **not** need an apply.

### ```colonize clean``` ###

**clean** is akin to ```make clean``` and should remove all of the generated
files that are created in the ```prep``` step.  This happens regardless if
a destroy or apply were done before hand.

### ```colonize generate <resource> [options]```

The **generate** command is used to provide convienience to generating Colonize
resources and project structures. **Generate** provides sub-commands for each
`resource-type` to create.


#### ```colonize generate branch <name> [options]```

The **branch** generation sub-command is used to generate a Colonize branch, including
build order file, environment directory, environment tfvars, and optionally a list
of leafs underneath the branch.

The following command:

`$ colonize generate branch myapp --leafs security_groups,database,instances`

Will generate the following branch & leaf structures
```
myapp
├── build_order.txt
├── database
│   └── main.tf
├── env
│   ├── dev.tfvars
│   ├── test.tfvars
│   └── prod.tfvars
├── instances
│   └── main.tf
└── security_groups
    └── main.tf
```

#### ```colonize generate leaf <name>```

The branch generation sub-command is used to generate a Colonize branch, including
build order file, environment directory, environment tfvars, and optionally a list
of leafs underneath the branch. When using the `generate branch` command with the
`--leafs` option, this command is internally called for each leaf.

The following command

`$ colonize generate leaf myleaf`

Will generate the following structure in the `mybranch` branch
```
mybranch
├── build_order.txt
├── env
└── myleaf
    └── main.tf
```


## CHANGELOG ##

  * 0.0.0 - still in development.

## CONTRIBUTORS ##

  * A special thanks to [2ndwatch](http://2ndwatch.com/) for supporting this
project!
  * Craig Monson
  * Joey Yore
  * Lars Cromley
