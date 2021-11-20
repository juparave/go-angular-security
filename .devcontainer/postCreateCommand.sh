#!/bin/bash

# This script is run after the container is created.
# It is run as root in the container.

# Fix files ownership and permissions
cd /workspaces && chown -R $(whoami) . && setfacl -bnR .
find . -type d -exec echo -n '"{}" ' \; | xargs chmod 755
find . -type f ! -name "*.sh" -exec echo -n '"{}" ' \; | xargs chmod 644 

