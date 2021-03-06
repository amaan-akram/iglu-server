#!/bin/bash

## This scripts generates Markdown formatted documentation of all packages in
## the project, using 'go doc' and 'godoc2markdown' tools.

## Change these variables...
# This shall be the module name of the project.
module="github.com/Nacdlow/iglu-server"
# This shall be the path of the wiki site.
docpath="../nacdlow-server.wiki"

mkdir -p $docpath

# Generate documentation/markdown of project
for fullpkg in $(go list ./...)
do
	package=$(echo $fullpkg | rev | cut -d"/" -f1 | rev)
	path=$(echo $fullpkg | sed "s#$module##")
	mkdir -p $docpath/$path
	echo "Generating $path..."
	echo "[< Go back to docs listing](/home)" > $docpath/$path/$package.md
	go doc -all $fullpkg | godoc2markdown >> $docpath/$path/$package.md
done

# Generate home wiki page with navigation listing
docs=$(cd $docpath && find -L . | grep \\.md)
echo -e "# Welcome to the documentation wiki\n\n" > $docpath/home.md
echo -e "This documentation is generated for commit
$(git rev-parse --short HEAD) (*$(git log -1 --pretty=%B)*).  \n\n" >> $docpath/home.md
echo "### Documentation listing" >> $docpath/home.md
echo 
for item in $docs
do
	propername=$(echo $item | sed "s#./##" | sed "s#\\.md##")
	if [ $propername != "home" ]; then
		echo "- [$propername]($propername)" >> $docpath/home.md
	fi
done
echo -e "\n\n> Generated using [godoc2markdown](https://humaidq.ae/projects/godoc2markdown/)
and a script." >> $docpath/home.md

echo "Done. Wiki generated at $docpath"
