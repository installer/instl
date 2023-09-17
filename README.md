<h1 align="center">INSTL <br> The Easiest Installer for GitHub Projects</sup></h1>

<p align="center">

<a style="text-decoration: none" href="https://github.com/installer/instl/releases">
<img src="https://img.shields.io/badge/platform-windows%20%7C%20macos%20%7C%20linux-informational?style=for-the-badge" alt="Downloads">
</a>

<a style="text-decoration: none" href="https://instl.sh/stats">
<img src="https://img.shields.io/endpoint?url=https://instl.sh/api/v1/badge/shields.io/stats/total&style=for-the-badge" alt="Handled installations">
</a>

</p>

<p align="center">
Instl is an installation script generator for GitHub projects. <br/>
It does not need any setup, and can be used to install most GitHub projects on Linux, macOS and Windows. <br/>
You can easily add installation commands to your <code>README.md</code> - they just work!
</p>

<p align="center">
<img width="900" src="https://raw.githubusercontent.com/installer/instl/main/demo.gif" alt="Instl Demo">
</p>

----

<p align="center">
<strong><a href="https://docs.instl.sh">Documentation</a></strong>
|
<strong><a href="https://docs.instl.sh/contributing/writing-code">Contributing</a></strong>
</p>

----

## Key Features

- 💻 Cross-Platform: Works on Windows, macOS and Linux out of the box
- 🧸 One-Click Installs: Install any GitHub project with just a single command
- 🪶 Web-Based: No need to install Instl - scripts are generated server-side
- ⚙️ Intelligent Configuration: Instl adapts to the project's structure
- 🕊️ On the fly: Installer scripts are created in real-time for each project
- 📊 [Track Your Installs](https://instl.sh/stats): Installation metrics for your projects at your fingertips

## Try it out:

Install our demo repository, `instl-demo`, to see instl in action. If successful, you should be able to run `instl-demo` right from your terminal.

| Platform | Command                                                    |
| -------- |------------------------------------------------------------|
| Windows  | <code>iwr instl.sh/installer/instl-demo/windows</code>     | iex |
| macOS    | <code>curl -sSL instl.sh/installer/instl-demo/macos</code> | bash |
| Linux    | <code>curl -sSL instl.sh/installer/instl-demo/linux</code> | bash |


## Usage

The fastest way to create your own instl commands, is by visiting [instl.sh](https://instl.sh) and using the builder.

Alternatively, you can create your own commands by using the following URL structure:

> [!NOTE]
> Replace `{username}` and `{reponame}` with your GitHub username and repository name.


#### Windows

```powershell
iwr instl.sh/{username}/{reponame}/windows | iex
```

#### macOS

```bash
curl -sSL instl.sh/{username}/{reponame}/macos | bash
```

#### Linux

```bash
curl -sSL instl.sh/{username}/{reponame}/linux | bash
```
