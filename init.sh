#!/bin/bash
read -p "Insert the new name of the module:" modulename

echo

read -p "The name of you module will be $modulename is this corect?(y/n) " yn
echo "Rename all references!"
case $yn in
y) find . -type f -not -path "./.git/*" -not -path "./init.sh" -not -path "./init.py" -type f -print0 | xargs -0 sed -i 's/go-printos-backend-quickstart/'"$modulename"'/g' ;;
n)
    echo exiting...
    exit
    ;;
*)
    echo invalid response
    exit 1
    ;;
esac
echo "Enabling git hooks"
git config core.hooks .hooks