Next Gen Particle CLI
==========

[Based on the Heroku CLI](https://github.com/heroku/heroku-cli)

**This is not an official tool by Particle. This is in alpha phase**

Overview
========

*The Next Gen Particle CLI is a delicacy consisting of a light Go shell filled with a sweet creamy Node.js core.*

The Go shell manages its own Node.js installation. It is able to update itself and the Node plugins. It forwards commands to the Node.js core. Since it is compiled natively for each platform it is easy to install.

The Node.js core is today's production Particle CLI. It performs all interactions with the Particle cloud.

Overview
========

heroku-cli (also known as Toolbelt v4) is the next generation Heroku CLI. Currently, it is not feature complete with the existing CLI and not intended to be used standalone. It is used internally inside the existing Toolbelt v3 Ruby CLI for some commands.

The goals of this project were to make plugins more flexible, remove Ruby as a runtime dependency, and make the CLI faster.

We are in the process of migrating commands over, some you can try out include: `heroku status`, `heroku fork`, and `heroku git:remote`.

Architecture
============

The v4 CLI is written both in Go and node. The core part of the CLI is this Go project and manages a node installation in `~/.heroku`. It handles taking in arguments, updating itself and node plugins, as well as some core commands like `heroku login`.

Most commands are plugins. An example of one would be [heroku-git](https://github.com/heroku/heroku-git). You can install it with `heroku plugins:install` and use it with `heroku git:remote`. When you install it, it uses npm to install the plugin into `~/.heroku/node_plugins/heroku-git`. Then, the commands become available to the v4 CLI.

We have shims in place for some commands inside the current v3 Toolbelt which will setup v4, install the plugin, and exec the command. [Here is the shim for heroku-git](https://github.com/heroku/heroku/blob/master/lib/heroku/command/git.rb#L47).

Once we have plugins to replicate every command in the v3 Toolbelt, we will then ship this CLI standalone.

For more on developing plugins, [read the blog announcement](https://blog.heroku.com/archives/2015/3/19/making_cli_plugins_better_with_a_new_architecture)

Issues
======

For problems directly related to the CLI, [add an issue on GitHub](https://github.com/heroku/heroku-cli/issues/new).

For other issues, [submit a support ticket](https://help.heroku.com/).

[Contributors](https://github.com/heroku/heroku-cli/contributors)
