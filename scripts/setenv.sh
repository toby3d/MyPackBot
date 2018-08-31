#!/bin/bash
#
# This script just set GitLab CI variables for correct work of Makefile
# commands. Run it by "$ . ./sripts/setenv.sh".

echo "Set variables..."

export CI_PROJECT_NAMESPACE "toby3d"
export CI_PROJECT_NAME "mypackbot"

echo "Variables is set!"