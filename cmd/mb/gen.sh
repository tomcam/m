#!/bin/zsh

# Do this only in a specified directory
# because it generates files

DIR=~/code/m/cmd/mb/theme-test
mkdir -p ./$DIR
INDEX=$DIR/index.md

# Generate home page for tests
cat <<-EOF > $INDEX
---
Theme: pillar
Sidebar: none
Mode: light
---
# Themes being test

EOF

# Create an array of filenames
declare -a themes=(
	'wide'
	'pillar'
	)


# Loop through the array
for file in "${themes[@]}"
do
	# Generate a filename base on traits
  # to test
  FILENAME=theme-${file}-left-light
cat <<-EOM > $DIR/$FILENAME.md
---
Theme: ${file}
Sidebar: left
Mode: light
Hey: there
---
# Auto-generated test for ${file} theme

Here's what we know about the test.
* Theme is **${file}** Acording to frontmatter it's **{{ .Page.FrontMatter.Theme }}**
* Sidebar should be on the **left**.
* Mode should be **light**

**Note**

* Resize this to check for responsive versions


```
print "hello, world."
```

EOM
  mv $FILENAME $DIR
  LINK="* **Theme:** ${file} [${file}-left-light](${FILENAME}.html)"
  echo "About to append ${LINK} to  ${INDEX}"
  echo $LINK >> $INDEX
done


#echo "Files created in ${DIR}:"
#ls $DIR
mb build $DIR
open $DIR/.mb/pub/index.html
#nvim $INDEX
exit

read -r -d '' VAR <<-'EOF'
  abc'asdf"
  $(dont-execute-this)
  foo"bar"''
EOF

