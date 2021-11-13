#!/bin/zsh

# Do this only in a specified directory
# because it generates files

DIR=~/code/m/cmd/mb/theme-test
mkdir -p ./$DIR

cat <<-EOF > ./$DIR/index.md
---
Theme: pillar
Sidebar: none
Mode: light
---
EOF

# Create an array of filenames
declare -a files=(
	'wide'
	'pillar'
	)


# Loop through the array
for file in "${files[@]}"
do
	# And execute the command
  FILENAME=theme-${file}-left-light.md
cat <<-EOF > ./$DIR/$FILENAME
---
Theme: ${file}
Sidebar: left
Mode: light
---
# Auto-generated test for ${file} theme

Here's what we know about the test.
* Theme is  ${file}, and according to frontmatter it's {{ Page.FrontMatter.Theme }}
* Sidebar should be on the left.
* Mode should be light

```
# Code test
print "hello, world."
```

EOF
  echo "${FILENAME}"
done
ls foo*.md
exit

read -r -d '' VAR <<-'EOF'
  abc'asdf"
  $(dont-execute-this)
  foo"bar"''
EOF

