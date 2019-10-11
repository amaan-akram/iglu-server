#!/bin/sh

# Change these variables...
module="gitlab.com/group-nacdlow/source"
docpath="../source.wiki"

mkdir -p $docpath
for fullpkg in $(go list ./...)
do
	package=$(echo $fullpkg | rev | cut -d"/" -f1 | rev)
	path=$(echo $fullpkg | sed "s#$module##")
	mkdir -p $docpath/$path
	echo "Generating $path..."
	go doc -all $fullpkg | godoc2markdown > $docpath/$path/$package.md
done

# Create home wiki page
docs=$(cd $docpath && find -L . | grep \\.md)
echo "# Welcome to the documentation wiki\n\n" > $docpath/home.md
echo "This documentation is generated at $(date) for commit
$(git rev-parse --short HEAD) (*$(git log -1 --pretty=%B)*).  \n\n" >> $docpath/home.md
echo "### Documentation listing" >> $docpath/home.md
echo 
for item in $docs
do
	propername=$(echo $item | sed "s#./##" | sed "s#\\.md##")
	echo "- [$propername]($propername)" >> $docpath/home.md
done
echo "\n\n> Generated using [godoc2markdown](https://humaidq.ae/projects/godoc2markdown/)
and a script." >> $docpath/home.md

echo "Done. Wiki generated at $docpath"
