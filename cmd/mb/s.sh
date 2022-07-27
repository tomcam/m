#!/bin/zsh
# Converts a text file with raw text to convert to Markdown
echo "s.sh Needs better documentation"
# s.sh
# -e means do one pass for each command
sed -e 's/```\n\([[:print:]]*\)\n```\n/```\n HI ```\n/g' \
  -e 's/{{ \([[:print:]]*\) }}/{{"{{"}} \1 {{"}}"}}/g' \
  -e 's/{{- \([[:print:]]*\) }}/{{"{{-"}} \1 {{"}}"}}/g' \
  -e 's/{{ \([[:print:]]*\) -}}/{{"{{"}} \1 {{"-}}"}}/g' \
  -e 's/{{- \([[:print:]]*\) -}}/{{"{{-"}} \1 {{"-}}"}}/g'    < $1  | pbcopy
