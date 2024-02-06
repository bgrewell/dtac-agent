# Overview

The Distributed Telemetry and Advanced Control (DTAC) framework is a collection of projects designed to reduce the time
to completion of software projects and testbeds by providing a highly reusable and extensible framework for the
collection of monitoring and manipulation of a wide variety of systems.

This project, the DTAC Agent, is focused on the endpoints. It is designed to run on various operating systems including
Windows, Linux and MacOS (Darwin). It provides access to a wide variety of telemetry on these systems and also provides
the ability to control many operating system and application parameters out of the box. The goal of the DTAC agent and
broader DTAC ecosystem is to replace legacy insecure APIs, custom tooling and shell scripts that do things like `ssh user@host <command>`
with a more feature complete, flexible and secure framework for automation and telemetry. To this end it has been designed
to be highly extensible through a multitude of methodologies described in more detail in the [extensibility](#extensibility)
section below.

Out of the box DTAC Agent supports both gRPC and REST as the Frontend API Protocols and has support for easily adding
additional frontend protocols.

## Installation

<tabs>
<tab title="Linux">
    <procedure title="Debian Based Linux Installation" id="debian-based-installation">
        <step>
            <p>Download the latest <code>.deb package</code> from the releases section of the GitHub repository.</p>
        </step>
        <step>
            <p>Open a terminal and navigate to the directory where the <code>.deb package</code> was downloaded. Typically <code>~/Downloads</code> on desktop versions of Linux.</p>
        </step>
        <step>
            <p>Run the following commands to install the package</p>
            <code-block lang="console" prompt="$">
                sudo apt install ./dtac-agentd_1.2.0_linux_amd64.deb
            </code-block>
            <note>
                <p>
                    Replace package name with downloaded version
                </p>
            </note>
        </step>
    </procedure>
    <procedure title="Redhat Based Linux Installation" id="redhat-based-installation">
        <step>
            <p>Download the latest <code>.rpm package</code> from the releases section of the GitHub repository.</p>
        </step>
        <step>
            <p>Open a terminal and navigate to the directory where the <code>.rpm package</code> was downloaded.</p>
        </step>
        <step>
            <p>Run the following commands to install the package</p>
            <code-block lang="console" prompt="$">
                sudo rpm -i ./dtac-agentd_1.2.0_linux_amd64.rpm
            </code-block>
            <note>
                <p>
                    Replace package name with downloaded version
                </p>
            </note>
        </step>
    </procedure>
</tab>
<tab title="Windows">
    <procedure title="Windows Installation" id="windows-based-installation">
    </procedure>
</tab>
<tab title="MacOS">
    <procedure title="MacOS Installation" id="macos-based-installation">
    </procedure>
</tab>
</tabs>

## Building From Source

<note>
    <p>TODO</p>
</note>

## Glossary

A definition list or a glossary:

First Term
: This is the definition of the first term.

Second Term
: This is the definition of the second term.
