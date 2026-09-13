#!/bin/bash


YAML_FILE="$1".yaml
PNG_FILE="$1".png

if [ -f "$PNG_FILE" ]; then
   rm "$PNG_FILE"
fi

../yamlfix/bin/yamlfix "$YAML_FILE"
awsdac "$YAML_FILE" -o "$PNG_FILE" --force

if [ -f "$PNG_FILE" ]; then
   xdg-open "$PNG_FILE"
fi
