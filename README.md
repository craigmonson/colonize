# Colonize #

Colonize is a configurable, albeit opinionated way to organize and manage
your terraform templates.  It revolves around the idea of environments, and
allows you to organize templates, and template data around that common idiom.

Once it's been configured, it allows for hierarical templates and variables, and
the ability to organize them in a defined manageable way.

## Install ##

### setup via go ###
### add .colonize.yaml to rootof project ###

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

# OLD DOCUMENTATION BELOW #
#### 4.  There must be a make\_order.txt file in order to use the branch make ####

## Execution ##

The terraform build scripts utilize the ```make``` command to perform all your
terraform tasks.  Simply typing ```make``` or ```make help``` will print out a
simple help message.  It uses some of the same subcommands found in terraform:
**plan** **apply** and **destroy**, adds the **prep** goal, and adds the
familiar **clean** goal as well.

**Note** that the following commands work the same in both the leaves and
branches.

#### %> make prep environemnt=\<environment\> ####
The **prep** make command is a prerequisite for the **plan** **apply**, and
**destroy** commands.  It is used to setup the terraform run based on the
\<environment\> that's passed in (see **plan** below for more info on that).
This command isn't necessary, as both **plan** and **destroy** will execute
this step on their own.  It is useful for debugging, as this command generates
all the build files.

All of the files generated by the build system are prevfixed with an underscore,
to denote them by convention as being constructed via make.  For example,
running ```%> make prep environment=dev``` in the ```app1/elb``` directory
will generate these files (this list is subject to change):

```bash
_base_variable_setup.tf
_combined.tfvars
_combined_templates.tf
_iam_remote_state.tf
_sg_remote_state.tf
_terraform_remote_config.sh*
_variable_setup.tf
_vpc_remote_state.tf
```

#### %> make plan environment=\<environment\> ####

The **plan** make command requires that you pass in the environment that you
wish to run this template in.  The \<environment\> is the name of the
```.tfstate``` file that's in our root, but without the ```.tfstate```.  Using
the tree from above, we can plan the app1/db/instance:

```bash
# pwd == ./
# setup keys for aws
%> export AWS_ACCESS_KEY_ID=NEVERGONNARUNAROUNDANDDESERTYOU
%> export AWS_SECRET_ACCESS_KEY=NEVERGONNAMAKEYOUCRY
# must be in the app1/db/instance dir to execute
%> cd app1/db/instance
%> make plan environment=prod
../terraform_build/remote_config.sh prod
Creating (or overwriting) temp_setup.tf with provider and variables
configuring to use state: foundation_prod/puppet_master_prod.tfstate
Initialized blank state with remote state enabled!
Remote state configured and pulled.
terraform get -update
...
...
Plan: 2 to add, 0 to change, 0 to destroy.
%>
```

You'll note that there are a bunch of files that make has generated in order
to automatically set some things up for the template, namely the
```terraform.tfplan``` file, which we'll need in the next step.

**Note** that you can skip the remote setup by setting **skip_remote=true**.
This will skip running the remote configuration command.  For example: you'd
need to do this when spinning up a new vpc, since the remote bucket is setup by
the vpc, you've got a chicken and egg problem.  You could do this as well if
you are testing some terraform, and don't want to clutter up the remote bucket.

#### %> make apply ####

This simply applies the plan which was stored in ```terraform.tfplan``` for you,
so you can be on your merry way.  It'll also store the state in s3 for you as
well.  You **MUST plan before you apply.**

**Note** that you can use the **skip_remote=true** variable here too and it'll
do the same as explained above.  However, you can also set
**remote_state_after_apply=true** to configure the remote state *AFTER* the
apply is complete.  This will sync the remote state.  Again, this is useful for
any new vpc's, as the bucket won't be available till after the apply is done.

```bash
%> make apply
terraform apply terraform.tfplan
terraform_remote_state.vpc: Creating...
...
...
Apply complete! Resources: 2 added, 0 changed, 0 destroyed.
...
...
%>
```

#### %> make destroy environment=\<environment\> (validate\_destroy=true) ####

This will run the ```terraform destroy``` command, effectively destroying all
the resources this leaf or branch controls.  It will prompt you to ensure you
want to destroy the current template by default.  If you wish to skip the
prompted validation, then pass in the ```validate_destroy=true``` argument when
running ```make destroy```

If there is a ```pre_destroy.sh``` file in the current directory, ```make
destroy``` will execute that script *before* the ```terraform destroy```
command is run.

If there is a ```post_destroy.sh``` file in the current directory, ```make
destroy``` will execute that script *after* the ```terraform destroy``` command
is run.

```bash
%> make destroy environment=prod
Refreshing Terraform state prior to plan...

terraform_remote_state.vpc: Refreshing state...
...
...
aws_instance.single: Destruction complete
terraform_remote_state.vpc: Destroying...
terraform_remote_state.vpc: Destruction complete
%>
```

#### %> make clean ####

We should clean up after ourselves no?  This simply removes all the files used
to execute the template.  It is recommended that you run clean often, ensuring
that old build files are cleaned up, and eliminating chances of environment
configuration pollution.  I can't stress this enough, **CLEAN OFTEN**

```bash
%> make clean
rm -f terraform.tfplan                                                  
rm -f destroy.tfplan                                                    
rm -rf .terraform                                                       
rm -f _temp_setup.tf                                                    
rm -f _vpc_remote_state.tf                                              
rm -f _combined.tfvars
...
%>
```

## A Note about statefiles ##

This build creates the ```terraform_remote_config.sh``` which sets up the remote
state configuration for a particular leaf.  In v1, the states were stored in
an s3 bucket named ```terraform_states_<app>_<vpc_env>``` and the key was
```<environment>/<template>_<environment>.tfstate```

In v2, the bucket remains the same, but the keys change to match the structure
of the tree.  The new key is structured as:
```<environment>/<project_tree>/<template>_<environment>.tfstate```

So, for example, the ```app1/db/iam``` template:

```bash
...
├── app1
│   ├── db
│   │   ├── iam
│   │   │   └── ...
...
```

would generate the s3 key: ```dev/app1/db/iam_dev.tfstate```, in the bucket
```terraform_states_foobar_dev```

**NOTE** that converting from v1 to v2 will require some manual handling of the
statefiles.
